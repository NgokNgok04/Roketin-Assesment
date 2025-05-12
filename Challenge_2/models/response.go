package models

type PaginatedMoviesResponse struct {
	Data []Movie `json:"data"`
	Meta struct {
		Page  int   `json:"page"`
		Limit int   `json:"limit"`
		Total int64 `json:"total"`
	} `json:"meta"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}