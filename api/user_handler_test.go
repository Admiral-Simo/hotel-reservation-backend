package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	if err != nil {
		t.Error(err)
	}
    defer resp.Body.Close()

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Fatal("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Fatal("expected the EncryptedPassword not to be included in the json response")
	}

	user.EncryptedPassword = ""
	if reflect.DeepEqual(user, params) {
		t.Fatalf("expected user '%v' but got '%v'", params, user)
	}
}

func init() {
	// load envirement variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("couldn't load envirement variables")
	}

}
