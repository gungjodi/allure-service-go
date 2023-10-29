package models

type BatchDeleteRequest struct {
	ProjectIds []string `json:"projects" example:"test123,test456,test789"`
	Async      bool     `json:"async"`
}

type CreateProjectRequest struct {
	ID string `json:"id" example:"default"`
}
