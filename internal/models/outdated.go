package models

type OutDatedResponse struct {
	RepoURL            string               `json:"repo_url,omitempty"`
	OutdatedDependency []OutdatedDependency `json:"outdated_dependency,omitempty"`
}

type OutdatedDependency struct {
	DependencyName              string `json:"dependency_name"`
	RepositoryDependencyVersion string `json:"repository_dependency_version"`
	RegistiryDependencyVersion  string `json:"registiry_dependency_version"`
}

type OutDatedData struct {
	ScanID               string           `json:"scan_id"`
	Email                []string         `json:"email"`
	OutdatedDependencies OutDatedResponse `json:"outdated_dependencies,omitempty"`
}
