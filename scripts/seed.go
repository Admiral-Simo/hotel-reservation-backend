package main

import (
	"context"
	"fmt"
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
	ctx        = context.Background()
	counter    int
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
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.DeluxRoomType,
			BasePrice: 1999.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 299.9,
		},
	}

	fmt.Printf("hotel number %d: %+v\n", counter, insertedHotel)

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("inserted room:", insertedRoom)
	}

	counter++
}

func main() {
	seedHotel("Royal Mansour", "Marrakech Morocco", 3)
	seedHotel("Mazagan Beach Resort", "Casablanca", 4)
	seedHotel("Dont die in your seleep", "London", 1.5)
}

func init() {
	counter = 1
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

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
