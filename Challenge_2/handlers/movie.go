package handlers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/models"
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
		title := c.FormValue("title")
		description := c.FormValue("description")
		durationStr := c.FormValue("duration")
		artistsNames := c.FormValue("artists")
		genreNames := c.FormValue("genres")
		if title == "" {return utils.HandleClientError(c, "title cant be empty")}
		if description == "" {return utils.HandleClientError(c, "description cant be empty")}
		if durationStr == "" {return utils.HandleClientError(c, "duration cant be empty")}
		
		duration, err := strconv.Atoi(durationStr)
		if err != nil || duration < 1 {return utils.HandleClientError(c, "invalid duration value")}
		
		file, err := c.FormFile("video")
		if err != nil {return utils.HandleClientError(c, "video file is required")}

		if err := utils.ValidateVideo(file); err != nil {return utils.HandleClientError(c, err.Error())}

		if err := db.Where("title = ?", title).First(&models.Movie{}).Error;  err == nil {
			return utils.HandleClientError(c, "movie with that title already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleError(c, "database error")
		}		

		artistList := strings.Split(artistsNames,",")
		genreList := strings.Split(genreNames,",")

		var artists []models.Artist
		if err := FindOrCreateArtist(db, artistList, &artists); err != nil {
			return utils.HandleError(c, err.Error())
		}

		var genres []models.Genre
		if err := FindOrCreateGenre(db, genreList, &genres); err != nil {
			return utils.HandleError(c, err.Error())
		}

		os.MkdirAll("Challenge_2/store", os.ModePerm)
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		savePath := fmt.Sprintf("Challenge_2/store/%s",uniqueName)
		if err := c.SaveFile(file, savePath); err != nil {return utils.HandleError(c, "Failed to save video file")}

		movie := models.Movie {
			Title : title,
			Description: description,
			Duration: uint32(duration),
			VideoURL: savePath,
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
		if err != nil {return utils.HandleClientError(c, "Invalid movie ID")}
		
		title := c.FormValue("title")
		description := c.FormValue("description")
		durationStr := c.FormValue("duration")
		artistsNames := c.FormValue("artists")
		genreNames := c.FormValue("genres")

		var duration *uint32
		if durationStr != "" {
			minutes, err := strconv.Atoi(durationStr)
			if err != nil || minutes < 1 {
				return utils.HandleClientError(c, "Invalid duration")
			}
			tmp := uint32(minutes)
			duration = &tmp
		}

		var movie models.Movie
		if err := db.Preload("Artists").Preload("Genres").First(&movie, id).Error; err != nil {
			return utils.HandleError(c, "Movie not found")
		}

		if title != "" {
			var existingMovie models.Movie
			if err := db.Where("title = ?", title).First(&existingMovie).Error; err == nil && existingMovie.ID != uint32(id) {
				return utils.HandleClientError(c, "Movie with that title already exists")
			}
			movie.Title = title
		}
		if description != "" {movie.Description = description}
		if duration != nil {movie.Duration = *duration}

		file, err := c.FormFile("video")
		if err == nil {
			if err := utils.ValidateVideo(file); err != nil {return utils.HandleClientError(c, err.Error())}
			
			os.MkdirAll("Challenge_2/store", os.ModePerm)
			uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			savePath := fmt.Sprintf("Challenge_2/store/%s", uniqueName)

			if err := c.SaveFile(file, savePath); err != nil {return utils.HandleError(c, "Failed to save video file")}
			movie.VideoURL = savePath
		}

		if artistsNames != "" {
			artistList := strings.Split(artistsNames, ",")

			var artists []models.Artist
			if err := FindOrCreateArtist(db, artistList, &artists); err != nil {return utils.HandleError(c, err.Error())}
			movie.Artists = artists
		}
		
		if genreNames != "" {
			genreList := strings.Split(genreNames, ",")

			var genres []models.Genre
			if err := FindOrCreateGenre(db, genreList, &genres); err != nil {return utils.HandleError(c, err.Error())}
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

