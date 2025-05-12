package models

type Movie struct {
	ID          uint32   `json:"id" gorm:"primaryKey"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    uint32   `json:"duration"`
	VideoURL    string   `json:"video_url"`
	Artists     []Artist `json:"artists" gorm:"many2many:movie_artists;"`
	Genres      []Genre  `json:"genres" gorm:"many2many:movie_genres;"`
}