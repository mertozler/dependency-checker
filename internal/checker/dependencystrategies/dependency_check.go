package dependencystrategies

import (
	"errors"
	"fmt"
	"github.com/mertozler/internal/models"
)

type Checker interface {
	RegisterStrategy(name string, strategy DependencyStrategy)
	Check(strategyType string, url string) (*models.OutDatedResponse, error)
}

type Check struct {
	dependencyCheckStrategies map[string]DependencyStrategy
}

func NewCheckerStrategy() Check {
	return Check{dependencyCheckStrategies: make(map[string]DependencyStrategy)}
}

func (c *Check) RegisterStrategy(name string, strategy DependencyStrategy) {
	c.dependencyCheckStrategies[name] = strategy
}

func (c *Check) Check(strategyType string, url string) (*models.OutDatedResponse, error) {
	dependencyCheckStrategy := c.dependencyCheckStrategies[strategyType]

	if dependencyCheckStrategy == nil {
		return nil, errors.New("No strategy found")
	}
	outDatedData, err := dependencyCheckStrategy.Check(url)
	if err != nil {
		return nil, err
	}

	fmt.Println("Dependency check completed")
	return outDatedData, nil

}
