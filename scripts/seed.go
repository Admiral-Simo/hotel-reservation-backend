package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/db/fixtures"
	"github.com/Admiral-Simo/HotelReserver/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
	store        *db.Store
)

func main() {
    var (

        user = fixtures.AddUser(store, "foo", "bar", false)
        admin = fixtures.AddUser(store, "admin", "admin", true)
        hotel = fixtures.AddHotel(store, "tazrkount", "marrakesh ", 3.2, nil)
        room = fixtures.AddRoom(store, "medium", true, 199.99, hotel.ID)
        from = time.Now().AddDate(0, 0, 1)
        till = from.AddDate(0, 0, 4)
        booking = fixtures.AddBooking(store, user.ID, room.ID, from, till, 2)
    )
    fmt.Println("admin -> ", types.CreateTokenFromUser(admin))
    fmt.Println("james -> ", types.CreateTokenFromUser(user))
	fmt.Println("booking ->", booking.ID)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal("couldn't drop the database:", err)
	}
	err = client.Database(db.DBNAME).Drop(ctx)
	if err != nil {
		log.Fatalf("failed to drop %s: %s", db.DBNAME, err)
	}
	userStore = db.NewMongoUserStore(client)
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	bookingStore = db.NewMongoBookStore(client)
	store = &db.Store{
		User:    userStore,
		Booking: bookingStore,
		Room:    roomStore,
		Hotel:   hotelStore,
	}
}
