package utils

import "github.com/gofiber/fiber/v2"

func HandleError(c *fiber.Ctx, message string) error {
	return c.Status(500).JSON(fiber.Map{
		"error": message,
	})
}
func HandleClientError(c *fiber.Ctx, message string) error {
	return c.Status(400).JSON(fiber.Map{
		"error": message,
	})
}
