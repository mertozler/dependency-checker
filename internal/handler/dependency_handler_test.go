package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/magiconair/properties/assert"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/mocks"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostScanRequestHandler_Should_Success(t *testing.T) {
	//given
	requestURL := "http://localhost:8080/api/v1/dependency"
	scanRequest := models.Request{
		RepoURL: "github.com/atolye15/checklist",
		Email:   []string{"meozler@gmail.com", "ahmet@gmail.com"},
	}
	mockRepo := mocks.NewRepo(t)
	mockDependencyChecker := mocks.NewDependency(t)

	mockDependencyChecker.On("DependencyCheck", scanRequest.RepoURL).Return(&models.OutDatedResponse{}, nil)
	mockRepo.On("SetScanData", mock.Anything, mock.Anything).Return(nil)

	dependencyHandler := NewDependencyHandler(mockRepo, mockDependencyChecker)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/dependency", dependencyHandler.PostScanRequestHandler())

	scanRequestJSON, _ := json.Marshal(scanRequest)
	req := httptest.NewRequest("POST", requestURL, bytes.NewBuffer(scanRequestJSON))
	req.Header.Set("Content-Type", "application/json")

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestPostScanRequestHandler_Should_Return_Error_When_Repository_Lost_The_Connection(t *testing.T) {
	//given
	requestURL := "http://localhost:8080/api/v1/dependency"
	scanRequest := models.Request{
		RepoURL: "github.com/atolye15/checklist",
		Email:   []string{"meozler@gmail.com", "ahmet@gmail.com"},
	}
	mockRepo := mocks.NewRepo(t)
	mockDependencyChecker := mocks.NewDependency(t)

	mockDependencyChecker.On("DependencyCheck", scanRequest.RepoURL).Return(&models.OutDatedResponse{}, nil)
	mockRepo.On("SetScanData", mock.Anything, mock.Anything).Return(errors.New("error while connecting to repository"))

	dependencyHandler := NewDependencyHandler(mockRepo, mockDependencyChecker)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/dependency", dependencyHandler.PostScanRequestHandler())

	scanRequestJSON, _ := json.Marshal(scanRequest)
	req := httptest.NewRequest("POST", requestURL, bytes.NewBuffer(scanRequestJSON))
	req.Header.Set("Content-Type", "application/json")

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestPostScanRequestHandler_Should_Return_Error_When_Dependency_Checker_Fails(t *testing.T) {
	//given
	requestURL := "http://localhost:8080/api/v1/dependency"
	scanRequest := models.Request{
		RepoURL: "github.com/atolye15/checklist",
		Email:   []string{"meozler@gmail.com", "ahmet@gmail.com"},
	}
	mockRepo := mocks.NewRepo(t)
	mockDependencyChecker := mocks.NewDependency(t)

	mockDependencyChecker.On("DependencyCheck", scanRequest.RepoURL).Return(&models.OutDatedResponse{}, errors.New("error while checking dependencies"))

	dependencyHandler := NewDependencyHandler(mockRepo, mockDependencyChecker)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Post("/dependency", dependencyHandler.PostScanRequestHandler())

	scanRequestJSON, _ := json.Marshal(scanRequest)
	req := httptest.NewRequest("POST", requestURL, bytes.NewBuffer(scanRequestJSON))
	req.Header.Set("Content-Type", "application/json")

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestGetScanRequestHandler_Should_Success(t *testing.T) {
	//given
	scanID := "xxxx"
	requestURL := "http://localhost:8080/api/v1/dependency/" + scanID
	mockRepo := mocks.NewRepo(t)
	mockRepo.On("GetScanData", "xxxx").Return(&models.OutDatedData{}, nil)

	dependencyHandler := NewDependencyHandler(mockRepo, nil)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Get("/dependency/:scanid", dependencyHandler.GetScanRequestHandler())
	req := httptest.NewRequest("GET", requestURL, nil)

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestGetScanRequestHandler_Should_Return_Error_When_Repository_Throws_Error(t *testing.T) {
	//given
	scanID := "xxxx"
	requestURL := "http://localhost:8080/api/v1/dependency/" + scanID
	mockRepo := mocks.NewRepo(t)
	mockRepo.On("GetScanData", "xxxx").Return(nil, errors.New("error while getting scan data"))

	dependencyHandler := NewDependencyHandler(mockRepo, nil)
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Get("/dependency/:scanid", dependencyHandler.GetScanRequestHandler())
	req := httptest.NewRequest("GET", requestURL, nil)

	//when
	response, _ := app.Test(req, -1)

	//then
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}
