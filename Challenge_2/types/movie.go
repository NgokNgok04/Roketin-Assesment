package types

type CreateMovieType struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    uint32 `json:"duration"`
	ArtistIDs   []uint `json:"artists"`
	GenreIDs    []uint `json:"genres"`
}

type UpdateMovieType struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Duration    *uint32 `json:"duration"`
	ArtistIDs   []uint  `json:"artist_ids"`
	GenreIDs    []uint  `json:"genre_ids"`
}