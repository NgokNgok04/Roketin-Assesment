package types

type CreatMovieType struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	ArtistIDs   []uint `json:"artist_ids"`
	GenreIDs    []uint `json:"genre_ids"`
}