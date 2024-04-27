package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Admiral-Simo/HotelReserver/db/fixtures"
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user           = fixtures.AddUser(tdb.Store, "another", "user", false)
		nonAuthUser    = fixtures.AddUser(tdb.Store, "simo", "levels", false)
		hotel          = fixtures.AddHotel(tdb.Store, "bar hotel", "marrakesh", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now().AddDate(0, 0, 1)
		till           = from.AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till, 2)
		app            = fiber.New()
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	// middelewares
	app.Use(JWTAuthentication(tdb.User))

	// adding route
	app.Get("/:id", bookingHandler.HandleGetBooking)

	// making request
	req := httptest.NewRequest(http.MethodGet, "/"+booking.ID.Hex(), nil)
	req.Header.Add("X-Api-Token", types.CreateTokenFromUser(user))

	// getting response
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal("couldn't perform request", err)
	}

	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal("error decoding response:", err)
	}

	// avoid conflict due to time inaccuracy
	bookingResp.FromDate = time.Now()
	bookingResp.TillDate = bookingResp.FromDate
	booking.FromDate = bookingResp.FromDate
	booking.TillDate = bookingResp.TillDate

	if !reflect.DeepEqual(bookingResp, booking) {
		t.Fatalf("got -> %v\nexpected -> %v\n", bookingResp, booking)
	}

	// non authorized user test
	req = httptest.NewRequest(http.MethodGet, "/"+booking.ID.Hex(), nil)
	req.Header.Add("X-Api-Token", types.CreateTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)

	if err != nil {
		t.Fatal("couldn't perform request", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %d status code got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser      = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user           = fixtures.AddUser(tdb.Store, "another", "user", false)
		hotel          = fixtures.AddHotel(tdb.Store, "bar hotel", "marrakesh", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now().AddDate(0, 0, 1)
		till           = from.AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till, 2)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	// middelewares
	app.Use(JWTAuthentication(tdb.User))
	app.Use(AdminAuth)

	// adding route
	app.Get("/", bookingHandler.HandleGetBookings)

	// making request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", types.CreateTokenFromUser(adminUser))

	// getting response
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal("failed to perform request", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status code -> %d, expected -> %d.", resp.StatusCode, http.StatusOK)
	}

	var bookings []*types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal("cannot decode data ->", err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}

	// avoid conflict due to time inaccuracy
	bookings[0].FromDate = time.Now()
	bookings[0].TillDate = bookings[0].FromDate
	booking.FromDate = bookings[0].FromDate
	booking.TillDate = bookings[0].TillDate

	if !reflect.DeepEqual(bookings[0], booking) {
		t.Fatalf("got -> %v\nexpected -> %v\n", bookings[0], booking)
	}

	// test non admin user
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", types.CreateTokenFromUser(user))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %d status code got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
