package models

type Artist struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Name   string  `json:"name"`
	Movies []Movie `gorm:"many2many:movie_artists"`
}