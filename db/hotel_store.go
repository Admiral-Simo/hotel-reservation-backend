package db

import (
	"context"
	"fmt"

	"github.com/Admiral-Simo/HotelReserver/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	hotelColl = "hotels"
)

type HotelStore interface {
	Dropper

	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, filter bson.M, update bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping hotell collection")
	return s.coll.Drop(ctx)
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	if hotel.Rooms == nil {
		hotel.Rooms = []primitive.ObjectID{}
	}
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
