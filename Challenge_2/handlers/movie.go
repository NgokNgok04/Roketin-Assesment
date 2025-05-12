package handlers

import (
	"errors"
	"strconv"
	"strings"

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

func GetPaginatedMovies(db *gorm.DB) fiber.Handler {
	return func(c * fiber.Ctx) error {
		pageStr := c.Query("page","1")
		limitStr := c.Query("limit", "10")
		
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)
		
		if page < 1 {page = 1}
		if limit < 1 {limit = 5}
		offset := (page - 1) * limit

		var movies []models.Movie
		var total int64

		db.Model(&models.Movie{}).Count(&total)
		if err := db.Preload("Artists").Preload("Genres").Limit(limit).Offset(offset).Find(&movies).Error; err != nil {
			return utils.HandleError(c, "Failed to fetch movies")
		}

		return c.JSON(fiber.Map{
			"data": movies,
			"meta": fiber.Map{
				"page": page,
				"limit": limit,
				"total": total,
			},
		})
	}
}

func CreateMovie(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input types.CreateMovieType
		if err := c.BodyParser(&input); err != nil {
			return utils.HandleError(c, "Failed to create movie")
		}

		if input.Title == "" {
			return utils.HandleClientError(c, "Title cant be empty")
		}
		if input.Description == "" {
			return utils.HandleClientError(c, "Description cant be empty")
		}
		if len(input.Artists) == 0 {
			return utils.HandleClientError(c, "Artists cant be empty")
		}
		if len(input.Genres) == 0 {
			return utils.HandleClientError(c, "Genres cant be empty")
		}
		if input.Duration < 1 {
			return utils.HandleClientError(c, "Duration cant be 0 or negative")
		}
		if err := db.Where("title = ?", input.Title).First(&models.Movie{}).Error;  err == nil {
			return utils.HandleClientError(c, "Movie with that title already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleError(c, "Database error")
		}

		var artists []models.Artist
		if err := FindOrCreateArtist(db, input.Artists, &artists); err != nil {
			return utils.HandleError(c, err.Error())
		}

		var genres []models.Genre
		if err := FindOrCreateGenre(db, input.Genres, &genres); err != nil {
			return utils.HandleError(c, err.Error())
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
			return utils.HandleClientError(c, "Invalid movie ID")
		}
		
		var input types.UpdateMovieType
		if err := c.BodyParser(&input); err != nil {
			return utils.HandleError(c, "Invalid request body")
		}
		
		if input.Description != nil && *input.Description == "" {return utils.HandleClientError(c, "Description cant be empty")}
		if input.Title != nil && *input.Title == "" {return utils.HandleClientError(c, "Title cant be empty")}
		if input.Duration != nil && *input.Duration < 1 {return utils.HandleClientError(c, "Duration cant be zero or negative")}
		if input.Artists != nil && len(input.Artists) == 0 {return utils.HandleClientError(c, "Artists cant be empty")}
		if input.Genres != nil && len(input.Genres) == 0 {return utils.HandleClientError(c, "Genres cant be empty")}
		
		var existingMovie models.Movie
		if err := db.Where("title = ?", input.Title).First(&existingMovie).Error; err == nil && existingMovie.ID != uint32(id) {
			return utils.HandleClientError(c, "Movie with that title already exists")
		}
		
		var movie models.Movie
		if err := db.Preload("Artists").Preload("Genres").First(&movie, id).Error; err != nil {
			return utils.HandleError(c, "Movie not found")
		}

		if input.Title != nil {movie.Title = *input.Title}
		if input.Description != nil {movie.Description = *input.Description}
		if input.Duration != nil {movie.Duration = *input.Duration}

		var artists []models.Artist
		if err := FindOrCreateArtist(db, input.Artists, &artists); err != nil {
			return utils.HandleError(c, err.Error())
		}
		if input.Artists != nil {movie.Artists = artists}
		
		var genres []models.Genre
		if err := FindOrCreateGenre(db, input.Genres, &genres); err != nil {
			return utils.HandleError(c, err.Error())
		}
		
		if input.Genres != nil {movie.Genres = genres}
		if err := db.Save(&movie).Error; err != nil {
			return utils.HandleError(c, "Failed to update movie")
		}
		return c.JSON(movie)
	}
}

func SearchMovies(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		q := strings.ToLower(c.Query("q"))
		if q == "" {
			return utils.HandleClientError(c, "Missing query param 'q'")
		}

		var movies []models.Movie

		err := db.Preload("Artists").Preload("Genres").
		Joins("LEFT JOIN movie_artists ON movie_artists.movie_id = movies.id").
		Joins("LEFT JOIN artists ON artists.id = movie_artists.artist_id").
		Joins("LEFT JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Joins("LEFT JOIN genres ON genres.id = movie_genres.genre_id").
		Where(`LOWER(movies.title) LIKE ? OR 
		       LOWER(movies.description) LIKE ? OR 
			   LOWER(artists.name) LIKE ? OR 
			   LOWER(genres.name) LIKE ?`,
			   "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%").
		Group("movies.id").
		Find(&movies).Error

		if err != nil {
			return utils.HandleError(c, "Failed to search movies")
		}

		return c.JSON(movies)
		
	}
}

