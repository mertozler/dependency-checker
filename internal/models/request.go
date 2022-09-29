package models

type Request struct {
	RepoURL string   `json:"url,omitempty"`
	Email   []string `json:"emails,omitempty"`
}
