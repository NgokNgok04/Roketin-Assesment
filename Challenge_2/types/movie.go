package types

type CreateMovieType struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    uint32   `json:"duration"`
	Artists     []string `json:"artists"`
	Genres      []string `json:"genres"`
}

type UpdateMovieType struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Duration    *uint32  `json:"duration"`
	Artists     []string `json:"artists"`
	Genres      []string `json:"genres"`
}