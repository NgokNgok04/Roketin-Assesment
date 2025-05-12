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
		for _, name := range input.Artists {
			var artist models.Artist
			if err := db.Where("name = ?", name).FirstOrCreate(&artist, models.Artist{Name: name}).Error; err != nil {
				return utils.HandleError(c, "Failed to create movie")
			}
			artists = append(artists, artist)
		}

		var genres []models.Genre
		for _, name := range input.Genres {
			var genre models.Genre
			if err := db.Where("name = ?", name).FirstOrCreate(&genre, models.Genre{Name: name}).Error; err != nil {
				return utils.HandleError(c, "Failed to create movie")
			}
			genres = append(genres, genre)
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

