package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Admiral-Simo/HotelReserver/db/fixtures"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of %d but got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var genResp genericResponse

	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected 'error' but got '%s'", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected 'invalid credentials' but got '%s'", genResp.Msg)
	}
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

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of %d but got %d", http.StatusOK, resp.StatusCode)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}

	if authResponse.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth response")
	}

	// Set the encrypted password to an empty string, because we do NOT return that in any
	// JSON response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		t.Fatalf("expected authResponse to be '%v' but got '%v'", insertedUser, authResponse.User)
	}
}

func init() {
	// load envirement variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("couldn't load envirement variables")
	}

}
