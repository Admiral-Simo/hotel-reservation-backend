package api

import (
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	users := []types.User{}
	u := types.User{
		ID:        "",
		FirstName: "Mohamed",
		LastName:  "Khalis",
	}
	users = append(users, u)
	users = append(users, u)
	users = append(users, u)
	return c.JSON(users)
}

func HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	u := types.User{
		ID: id,
	}
	return c.JSON(u)
}
