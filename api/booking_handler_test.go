package api

import (
	"testing"
	"time"

	"github.com/Admiral-Simo/HotelReserver/db/fixtures"
)

func TestGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "mohamed", "khalis", true)
	hotel := fixtures.AddHotel(tdb.Store, "bar hotel", "marrakesh", 4, nil)
	room := fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)
	from := time.Now().AddDate(0, 0, 1)
	till := time.Now().AddDate(0, 0, 2)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till, 2)

	t.Log(booking)
}
