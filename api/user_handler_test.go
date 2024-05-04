package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/user", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "James",
		LastName:  "Foo",
		Password:  "lkjasdlkfjalskdfjalkjsd",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Nil(t, err, "getting response error: %v", err)

	defer resp.Body.Close()

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	assert.NotEmpty(t, user.ID, "expecting a user id to be set")

	assert.Empty(t, user.EncryptedPassword, "expected the EncryptedPassword not to be included in the json response (for security purposes)")

	user.EncryptedPassword = ""

	assert.NotEqual(t, user, params, "user and params should not be equal")
}

func init() {
	// load envirement variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("couldn't load envirement variables")
	}

}
