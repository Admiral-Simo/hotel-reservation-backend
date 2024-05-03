package db

import "context"

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
)

// Dropper is going to be for every store available
type Dropper interface {
	Drop(context.Context) error
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
	Logs    LogsStore
}
