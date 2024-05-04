package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Admiral-Simo/HotelReserver/db/fixtures"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateAuthenticateWithWrongCredentials(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.AddUser(tdb.Store, "mohamed", "khalis", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "mohamed@khalis.com",
		Password: "wrong password",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be bad request")

	var genResp genericResponse

	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "error", genResp.Type, "gen response type should be error")

	assert.Equal(t, "invalid credentials", genResp.Msg, "gen response message should be `invalid credentials`")
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := fixtures.AddUser(tdb.Store, "mohamed", "khalis", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "mohamed@khalis.com",
		Password: "mohamed_khalis",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err, "error getting response from /auth: %v", err)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "status code should be okay")

	var authResponse AuthResponse

	err = json.NewDecoder(resp.Body).Decode(&authResponse)

	assert.NoError(t, err, "error decoding body at /auth: %v", err)

	assert.NotEmpty(t, authResponse.Token, "expected the JWT token to be present in the auth response")

	// Set the encrypted password to an empty string, because we do NOT return that in any
	// JSON response
	insertedUser.EncryptedPassword = ""

	assert.Equal(t, insertedUser, authResponse.User, "response and insertedUser should be equal")
}

func init() {
	// load envirement variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("couldn't load envirement variables")
	}

}
