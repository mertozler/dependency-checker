package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mertozler/internal/checker"
	"github.com/mertozler/internal/checker/dependencystrategies"
	"github.com/mertozler/internal/checker/dependencystrategies/strategies"
	"github.com/mertozler/internal/config"
	"github.com/mertozler/internal/handler"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/internal/worker"
	"github.com/mertozler/pkg/mail"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	configs, err := config.LoadConfig("./configs")
	if err != nil {
		logrus.Fatal("Error while getting configs", err)
	}

	repo := repository.NewRepository(configs.Redis)

	checkerStrategy := dependencystrategies.NewCheckerStrategy()
	checkerStrategy.RegisterStrategy("github",
		strategies.NewGithub(configs.Github.API, configs.Javascript.RegistryURL, configs.Dependency.Path))
	dependencyChecker := checker.NewDependencyChecker(&checkerStrategy)

	dependencyHandler := handler.NewDependencyHandler(repo, dependencyChecker)

	mailSender := mail.NewMail(configs.Mail)
	outDatedWorker := worker.NewWorker(repo, mailSender, configs.MailSender)

	go outDatedWorker.OutdatedWorker()
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/dependency", dependencyHandler.PostScanRequestHandler())
	api.Get("/dependency/:scanid", dependencyHandler.GetScanRequestHandler())

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		logrus.Infof("Received %s signal", <-c)
		app.Shutdown()
	}()

	err = app.Listen(configs.Host.Port)
	if err != nil {
		logrus.Fatal("Error while listening port", err)
	}
}
