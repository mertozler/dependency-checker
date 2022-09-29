package dependencystrategies

import "github.com/mertozler/internal/models"

type DependencyStrategy interface {
	Check(url string) (*models.OutDatedResponse, error)
}
