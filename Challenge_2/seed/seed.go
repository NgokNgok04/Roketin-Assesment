package seed

import (
	"log"

	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	err := db.AutoMigrate(&models.Movie{}, &models.Artist{}, &models.Genre{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	artists := []models.Artist{
		{Name: "Christopher Nolan"},
		{Name: "Hans Zimmer"},
	}

	genres := []models.Genre{
		{Name: "Sci-Fi"},
		{Name: "Drama"},
	}

	db.Create(&artists)
	db.Create(&genres);

	movie := models.Movie{
		Title: "Interstellar",
		Description: "Space exploration and time dilation",
		Duration: 169,
		Artists: artists,
		Genres: genres,
	}

	db.Create(&movie)
}