package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mertozler/internal/checker"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/internal/repository"
	"github.com/sirupsen/logrus"
	"net/http"
)

type DependencyHandler struct {
	repo              repository.Repo
	dependencyChecker checker.Dependency
}

func NewDependencyHandler(repo repository.Repo, dependencyChecker checker.Dependency) *DependencyHandler {
	return &DependencyHandler{repo: repo, dependencyChecker: dependencyChecker}
}

func (d *DependencyHandler) PostScanRequestHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		logrus.Info("New scan request received")
		var req models.Request
		if err := c.BodyParser(&req); err != nil {
			logrus.Error("Error while parsing body", err)
			return c.Status(http.StatusBadRequest).JSON(models.Response{
				Status:  fiber.StatusBadRequest,
				Message: "Error while parsing body",
				Data: &fiber.Map{
					"error_message": err.Error(),
				},
			})
		}
		scanID := uuid.New().String()
		outDatedResponse, err := d.dependencyChecker.DependencyCheck(req.RepoURL)
		if err != nil {
			logrus.Error("Error while checking dependency", err)
			return c.Status(http.StatusInternalServerError).JSON(models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error while checking dependency",
				Data: &fiber.Map{
					"error_message": err.Error(),
				},
			})
		}
		outDatedData := models.OutDatedData{
			ScanID:               scanID,
			Email:                req.Email,
			OutdatedDependencies: *outDatedResponse,
		}

		err = d.repo.SetScanData(scanID, outDatedData)
		if err != nil {
			logrus.Error("Error while parsing body", err)
			return c.Status(http.StatusInternalServerError).JSON(models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error inserting data into repository",
				Data: &fiber.Map{
					"error_message": err.Error(),
				},
			})
		}
		logrus.Info("New scan request ended")
		return c.Status(fiber.StatusOK).JSON(models.Response{
			Status:  fiber.StatusOK,
			Message: "New Scan request ended successfully",
			Data: &fiber.Map{
				"scan_results": outDatedData,
			},
		})
	}
}

func (d *DependencyHandler) GetScanRequestHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		scanID := c.AllParams()["scanid"]
		logrus.Infof("GetScan request received for scan_id: %v", scanID)
		scanData, err := d.repo.GetScanData(scanID)
		if err != nil {
			logrus.Error("Error getting scan data from database", err)
			return c.Status(http.StatusInternalServerError).JSON(models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error getting scan data from database",
				Data: &fiber.Map{
					"Error Message": err.Error(),
				},
			})
		}
		logrus.Infof("GetScan request ended for scan_id: %v", scanID)
		return c.Status(http.StatusOK).JSON(models.Response{
			Status:  fiber.StatusOK,
			Message: "Get scan result from database successfully",
			Data: &fiber.Map{
				"scan_results": scanData,
			},
		})
	}
}
