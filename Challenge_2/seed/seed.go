package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/models"
	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/types"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load();
	
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"),
	)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	if err := db.AutoMigrate(&models.Movie{}, &models.Artist{}, &models.Genre{}); err != nil {
		log.Fatal("Migration failed:", err)
	}
	
	db.Exec(`
	  TRUNCATE TABLE 
	    movie_artists, 
	    movie_genres, 
	    movies, 
	    artists, 
	    genres 
	  RESTART IDENTITY CASCADE
	`)

	file, err := os.Open("Challenge_2/seed/dummy.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	var movies []types.CreateMovieType;
	if err := json.NewDecoder(file).Decode(&movies); err != nil {
		log.Fatal("Failed to decode JSON:", err)
	}
	for _, movie := range movies {
		var artists []models.Artist
		var genres []models.Genre
		for _, name := range movie.Artists {
			var artist models.Artist
			if err := db.Where("name = ?", name).FirstOrCreate(&artist, models.Artist{Name: name}).Error; err != nil {
				log.Fatal("Failed to create artist :", name)
			}
			artists = append(artists, artist)
		}
		
		for _, name := range movie.Genres {
			var genre models.Genre
			if err := db.Where("name = ?", name).FirstOrCreate(&genre, models.Genre{Name: name}).Error; err != nil {
				log.Fatal("Failed to create genre :", name)
			}
			genres = append(genres, genre)
		}

		createMovie := models.Movie {
			Title: movie.Title,
			Description: movie.Description,
			Duration: movie.Duration,
			VideoURL: "Challenge_2/store/dummy.mp4",
			Artists: artists,
			Genres: genres,
		}

		if err := db.Create(&createMovie).Error; err != nil {
			log.Fatal("Failed to create movie:", err)
		}
	}

	fmt.Println("Movies seeded.")
}