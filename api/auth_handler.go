package api

import (
	"errors"
	"net/http"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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

type AuthResponse struct {
	User *types.User `json:"user"`
}

type genericResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponse{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

// A handler should only do:
//   - serialization of the incoming request (JSON)
//   - do some data fetching from db
//   - call some business logic
//   - return the data back to the user
func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var aParams AuthParams
	if err := c.BodyParser(&aParams); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), aParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, aParams.Password) {
		return invalidCredentials(c)
	}

	tokenString := types.CreateTokenFromUser(user)

	resp := AuthResponse{
		User: user,
	}

	// Set the access token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    tokenString,
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(resp)
}
