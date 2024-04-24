package main

import (
	"context"
	"log"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedHotel(name, location string, rating float32) {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		log.Fatal("couldn't insert hotel:", err)
	}

	rooms := []types.Room{
		{
			Type:      "small",
			SeaSide:   true,
			BasePrice: 99.9,
		},
		{
			Type:      "medium",
			BasePrice: 1999.9,
		},
		{
			Type:      "large",
			BasePrice: 299.9,
		},
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedUser(firstName, lastName, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  "supersecurepassword",
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(ctx, user)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedHotel("Royal Mansour", "Marrakech Morocco", 3)
	seedHotel("Mazagan Beach Resort", "Casablanca", 4)
	seedHotel("Dont die in your seleep", "London", 1.5)
	seedUser("Mohamed", "Khalis", "personalsimoypo@gmail.com")
	seedUser("Toufik", "Khalis", "toufikhalis@gmail.com")
	seedUser("Driss", "El Haskouri", "drisshaskouri@gmail.com")
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
}
