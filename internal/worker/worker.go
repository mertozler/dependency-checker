package worker

import (
	"encoding/json"
	"github.com/mertozler/internal/config"
	"github.com/mertozler/internal/repository"
	"github.com/mertozler/pkg/mail"
	"github.com/sirupsen/logrus"
	"time"
)

type Worker struct {
	repo       repository.Repo
	sender     *mail.Mail
	mailSender *config.MailSender
}

func NewWorker(repo repository.Repo, sender *mail.Mail, mailSender *config.MailSender) *Worker {
	return &Worker{repo: repo, sender: sender, mailSender: mailSender}
}

func (w *Worker) OutdatedWorker() {
	for {
		time.Sleep(time.Duration(w.mailSender.Hour) * time.Hour)
		keys, err := w.repo.GetAllKeys()
		if err != nil {
			logrus.Error("Error while getting keys: ", err)
		}

		for _, key := range keys {
			data, repoErr := w.repo.GetScanData(key)
			if repoErr != nil {
				logrus.Error("Error while getting keys: ", repoErr)
			}

			marshalledOutdatedDependency, marshallErr := json.Marshal(data.OutdatedDependencies)
			if marshallErr != nil {
				logrus.Error("Error while marshalling outdated dependency: ", marshallErr)
			}

			mailErr := w.sender.SendMail(data.Email, "Dependency analysis for: "+data.OutdatedDependencies.RepoURL, string(marshalledOutdatedDependency))
			if mailErr != nil {
				logrus.Error("Error while sending mail: ", mailErr)
			}
		}

	}
}
