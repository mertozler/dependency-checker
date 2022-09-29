package models

type RegistiryResponse struct {
	ID          string   `json:"_id"`
	Rev         string   `json:"_rev"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DistTags    DistTags `json:"dist-tags"`
}
type DistTags struct {
	Latest string `json:"latest"`
	Next   string `json:"next"`
}
