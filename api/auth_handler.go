package api

import (
	"fmt"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var aParams AuthParams
	if err := c.BodyParser(&aParams); err != nil {
		return err
	}

	fmt.Println(aParams)

	return nil
}
