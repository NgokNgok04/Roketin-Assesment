package models

type Movie struct {
	ID          uint32   `json:"id" gorm:"primaryKey"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    uint32   `json:"duration"`
	Artists     []Artist `json:"artists" gorm:"many2many:move_movie_artists;"`
	Genres      []Genre  `json:"genres" gorm:"many2many:move_movie_genres;"`
}