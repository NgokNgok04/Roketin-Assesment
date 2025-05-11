package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NgokNgok04/Roketin-Assesment/Challenge_2/seed"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Movie struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Duration int `json:"duration"`
	Artists []string `json:"artists"`
 	Genres []string `json:"genres"`
}
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
	seed.Seed(db)

	app := fiber.New();

	movies := []Movie{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "hello world"})
	})
	app.Post("/api/movies", func(c *fiber.Ctx) error {
		movie:= &Movie{};

		if err := c.BodyParser(movie); err != nil {
			return err;
		}

		if movie.Title == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
		}
		movie.ID = len(movies) + 1
		movies = append(movies, *movie);

		return c.Status(201).JSON(movie);
	})
	fmt.Println("Hello Worlds");
	log.Fatal(app.Listen(":4000"));
}