package handlers

import (
	"strconv"

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
		var input types.CreateMovieType
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

func UpdateMovie(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return utils.HandleError(c, "Invalid movie ID")
		}

		var req types.UpdateMovieType
		if err := c.BodyParser(&req); err != nil {
			return utils.HandleError(c, "Invalid request body")
		}

		var movie models.Movie
		if err := db.Preload("Artists").Preload("Genres").First(&movie, id).Error; err != nil {
			return utils.HandleError(c, "Movie not found")
		}

		if req.Title != nil {movie.Title = *req.Title}
		if req.Description != nil {movie.Description = *req.Description}
		if req.Duration != nil {movie.Duration = *req.Duration}

		if len(req.ArtistIDs) > 0 {
			var artists []models.Artist
			db.Find(&artists, req.ArtistIDs)
			movie.Artists = artists
		}
		if len(req.GenreIDs) > 0 {
			var genres []models.Genre
			db.Find(&genres, req.GenreIDs)
			movie.Genres = genres
		}

		if err := db.Save(&movie).Error; err != nil {
			return utils.HandleError(c, "Failed to update movie")
		}

		return c.JSON(movie)
	}
}