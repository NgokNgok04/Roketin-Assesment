package handlers

import (
	"errors"

	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/models"
	"gorm.io/gorm"
)

func FindOrCreateGenre(db *gorm.DB, genreName []string, genres *[]models.Genre) error {
	for _, name := range genreName {
		var genre models.Genre
		if err := db.Where("name = ?", name).FirstOrCreate(&genre, models.Genre{Name: name}).Error; err != nil {
			return errors.New("failed to create genres")
		}
		*genres = append(*genres, genre)
	}

	return nil
}