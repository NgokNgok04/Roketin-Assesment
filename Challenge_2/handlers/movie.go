package handlers

import (
	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/models"
	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/types"
	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllMovies(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		var movies []models.Movie
		if err := db.Find(&movies).Error; err != nil {
			return utils.HandleError(c, "Failed to fetch movies")
		}
		return c.JSON(movies)
	}
}

func CreateMovie(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input types.CreatMovieType
		if err := c.BodyParser(&input); err != nil {
			return utils.HandleError(c, "Failed to create movie")
		}

		var artists []models.Artist
		if len(input.ArtistIDs) > 0 {
			if err := db.Find(&artists, input.ArtistIDs).Error; err != nil {
				return utils.HandleError(c, "Failed to find artists")
			}
		}

		var genres []models.Genre
		if len(input.GenreIDs) > 0 {
			if err := db.Find(&genres, input.GenreIDs).Error; err != nil {
				return utils.HandleError(c, "Failed to find genres")
			}
		}

		movie := models.Movie {
			Title : input.Title,
			Description: input.Description,
			Duration: uint32(input.Duration),
			Artists: artists,
			Genres: genres,
		}

		if err := db.Create(&movie).Error; err != nil {
			return utils.HandleError(c, "Failed to create movie")
		}

		return c.Status(201).JSON(movie)
	}
}