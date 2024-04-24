package api

import (
	"errors"
	"fmt"

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
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

// A handler should only do:
//   - serialization of the incoming request (JSON)
//   - do some data fetching from db
//   - call some business logic
//   - return the data back the user
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

	if !types.IsValidPassword(user.EncryptedPassword, aParams.Password) {
		return fmt.Errorf("invalid credentials")
	}

	// TODO: store jwt token into the User HEADER
	resp := AuthResponse{
		User:  user,
		Token: types.CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}
