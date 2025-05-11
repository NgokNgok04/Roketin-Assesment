package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
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