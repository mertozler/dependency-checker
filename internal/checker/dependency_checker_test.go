package checker

import (
	"errors"
	"github.com/mertozler/internal/models"
	"github.com/mertozler/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDependencyChecker_Should_Success(t *testing.T) {
	//given
	checker := mocks.NewChecker(t)
	expectedOutputData := &models.OutDatedResponse{
		RepoURL:            "github.com/mertozler/hichat",
		OutdatedDependency: nil,
	}
	checker.On("Check", "github", "github.com/mertozler/hichat").Return(expectedOutputData, nil)
	dependencyChecker := NewDependencyChecker(checker)
	//when
	actualOutputData, err := dependencyChecker.DependencyCheck("github.com/mertozler/hichat")
	//then
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedOutputData, actualOutputData)
}

func TestDependencyChecker_Should_Return_Error_When_No_Strategy_Type_Found_With_This_URL(t *testing.T) {
	//given
	checker := mocks.NewChecker(t)
	checker.On("Check", "hatalibirurlgiriyorumkistratejibulamasın", "hatalibirurlgiriyorumkistratejibulamasın").
		Return(nil, errors.New("No strategy type found"))
	dependencyChecker := NewDependencyChecker(checker)
	//when
	_, err := dependencyChecker.DependencyCheck("hatalibirurlgiriyorumkistratejibulamasın")
	//then
	assert.Equal(t, "No strategy type found", err.Error())
}

func TestDependencyCheckerGetStrategyType_Should_Success(t *testing.T) {
	//given
	checker := mocks.NewChecker(t)
	dependencyChecker := NewDependencyChecker(checker)
	//when
	strategyType := dependencyChecker.getStrategyType("github.com/mertozler/hichat")
	//then
	assert.Equal(t, "github", strategyType)
}
