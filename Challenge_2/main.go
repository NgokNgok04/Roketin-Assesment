package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/handlers"
	"github.com/gofiber/fiber/v2"
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

	app := fiber.New();
	// Route
	app.Get("/movies", handlers.GetPaginatedMovies(db))
	app.Post("/movies", handlers.CreateMovie(db))
	app.Put("/movies/:id", handlers.UpdateMovie(db))
	app.Get("/movies/search", handlers.SearchMovies(db)) 
	log.Fatal(app.Listen(":3000"));
}