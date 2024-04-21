package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/foo", handleFoo)
	app.Listen(":8080")
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "working just fine something"})
}
