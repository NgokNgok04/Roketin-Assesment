package handlers

import (
	"errors"

	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/models"
	"gorm.io/gorm"
)

func FindOrCreateArtist(db *gorm.DB, artistName []string, artists *[]models.Artist) error {
	for _, name := range artistName {
		var artist models.Artist
		if err := db.Where("name = ?", name).FirstOrCreate(&artist, models.Artist{Name: name}).Error; err != nil {
			return errors.New("Failed to create artists")
		}
		*artists = append(*artists, artist)
	}

	return nil
}