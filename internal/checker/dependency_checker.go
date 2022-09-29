package checker

import (
	"github.com/mertozler/internal/checker/dependencystrategies"
	"github.com/mertozler/internal/models"
	"github.com/sirupsen/logrus"
	"strings"
)

type Dependency interface {
	DependencyCheck(url string) (*models.OutDatedResponse, error)
}

type DependencyChecker struct {
	Checker dependencystrategies.Checker
}

func NewDependencyChecker(Checker dependencystrategies.Checker) *DependencyChecker {
	return &DependencyChecker{
		Checker: Checker}
}

func (d *DependencyChecker) DependencyCheck(url string) (*models.OutDatedResponse, error) {

	strategyType := d.getStrategyType(url)
	outDatedData, err := d.Checker.Check(strategyType, url)
	if err != nil {
		return nil, err
	}
	return outDatedData, nil
}

func (d *DependencyChecker) getStrategyType(url string) string {
	logrus.Infof("Getting strategy type url for %s", url)
	var splitedURL = strings.Split(url, ".")
	strategyType := splitedURL[0]
	return strategyType
}
