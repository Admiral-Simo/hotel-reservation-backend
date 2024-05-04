package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Admiral-Simo/HotelReserver/db/fixtures"
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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

	assert.NoError(t, err, "couldn't perform request: %v", err)

	var bookingResp *types.Booking

	err = json.NewDecoder(resp.Body).Decode(&bookingResp)
	assert.NoError(t, err, "error decoding response: %v", err)

	// avoid conflict due to time inaccuracy
	bookingResp.FromDate = time.Now()
	bookingResp.TillDate = bookingResp.FromDate
	booking.FromDate = bookingResp.FromDate
	booking.TillDate = bookingResp.TillDate

	assert.Equal(t, booking, bookingResp, "they should be equal")

	// non authorized user test
	req = httptest.NewRequest(http.MethodGet, "/"+booking.ID.Hex(), nil)
	req.Header.Add("X-Api-Token", types.CreateTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)

	assert.NoError(t, err, "couldn't perform request: %v", err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "they should be equal")
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

	assert.NoError(t, err, "failed to perform request:", err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "response code should be OK")

	var bookings []*types.Booking

	err = json.NewDecoder(resp.Body).Decode(&bookings)

	assert.NoError(t, err, "cannot decode data ->: %v", err)

	assert.Equal(t, 1, len(bookings), "should be only one booking")

	// avoid conflict due to time inaccuracy
	bookings[0].FromDate = time.Now()
	bookings[0].TillDate = bookings[0].FromDate
	booking.FromDate = bookings[0].FromDate
	booking.TillDate = bookings[0].TillDate

	assert.Equal(t, booking, bookings[0], "booking should be equal")

	// test non admin user
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", types.CreateTokenFromUser(user))

	resp, err = app.Test(req)
	assert.NoError(t, err, "failed to perform request: %v", err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "status code should be unauthorized")
}
