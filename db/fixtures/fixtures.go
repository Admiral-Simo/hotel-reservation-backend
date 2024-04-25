package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, from, till time.Time, numPersons int) *types.Booking {
	bookingParams := &types.Booking{
		UserID:     uid,
		RoomID:     rid,
		NumPersons: numPersons,
		FromDate:   from,
		TillDate:   till,
	}

	booking, err := store.Booking.InsertBooking(context.TODO(), bookingParams)
	if err != nil {
		log.Fatal(err)
	}

	return booking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hid primitive.ObjectID) *types.Room {
	roomParams := &types.Room{
		Size:    size,
		SeaSide: ss,
		Price:   price,
		HotelID: hid,
	}
	room, err := store.Room.InsertRoom(context.TODO(), roomParams)
	if err != nil {
		log.Fatal(err)
	}
	return room
}

func AddHotel(store *db.Store, name, loc string, rating float32, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := &types.Hotel{
		Name:     name,
		Location: loc,
		Rating:   rating,
		Rooms:    roomIDS,
	}
	hotel, err := store.Hotel.InsertHotel(context.TODO(), hotel)
	if err != nil {
		log.Fatal(err)
	}
	return hotel
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		IsAdmin:   admin,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
