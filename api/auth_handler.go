package api

import (
	"errors"
	"fmt"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

	user, err := h.userStore.GetUserByEmail(c.Context(), aParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(aParams.Password)); err != nil {
		return fmt.Errorf("invalid credentials")
	}

    // TODO: store jwt token into the User HEADER
	fmt.Println("authenticated ->", user)

	return nil
}
