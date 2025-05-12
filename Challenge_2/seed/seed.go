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


	artistNames := []string{
		"Robert Downey Jr.", //1
		"Chris Evans", 
		"Scarlett Johansson",
		"Chris Hemsworth", 
		"Tom Holland", //5
		"Mark Ruffalo",
		"Elizabeth Olsen", 
		"Benedict Cumberbatch", 
		"Paul Rudd", 
		"Chris Pratt", //10
		"Chadwick Boseman",
	}

	genreNames := []string{
		"Action", //1
		"Adventure", 
		"Drama", 
		"Comedy", 
		"Fantasy", //5
		"Sci-Fi", 
		"Thriller", 
		"War", 
		"Teen", 
		"Mystery", //10
	}

	for _, name := range artistNames {
		db.FirstOrCreate(&models.Artist{}, models.Artist{Name: name})
	}
	for _, name := range genreNames {
		db.FirstOrCreate(&models.Genre{}, models.Genre{Name: name})
	}

	log.Println("Artists and genres seeded.")


	file, err := os.Open("Challenge_2/seed/dummy.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	var movies []types.CreateMovieType;
	if err := json.NewDecoder(file).Decode(&movies); err != nil {
		log.Fatal("Failed to decode JSON:", err)
	}


	for _,movie := range movies {
		var artists []models.Artist
		var genres []models.Genre

		if len(movie.ArtistIDs) > 0 {
			if err := db.Find(&artists, movie.ArtistIDs).Error; err != nil {
				log.Fatal("Failed to find artists:", err)
			}
		}

		if len(movie.GenreIDs) > 0 {
			if err := db.Find(&genres, movie.GenreIDs).Error; err != nil {
				log.Fatal("Failed to find genres:", err)
			}
		}

		createMovie := models.Movie{
			Title:       movie.Title,
			Description: movie.Description,
			Duration:    movie.Duration,
			Artists:     artists,
			Genres:      genres,
		}

		if err := db.Create(&createMovie).Error; err != nil {
			log.Fatal("Failed to create movie:", err)
		}
	}
	fmt.Println("Movies seeded.")
}