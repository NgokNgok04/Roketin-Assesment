package models

type Genre struct {
	ID     uint32  `json:"id" gorm:"primaryKey"`
	Name   string  `json:"name"`
	Movies []Movie `gorm:"many2many:movie_genres"`
}