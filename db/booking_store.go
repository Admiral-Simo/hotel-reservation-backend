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
	bookingColl = "bookings"
)

type BookingStore interface {
	Dropper

	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingById(ctx context.Context, id string) (*types.Booking, error)
	IsRoomAvailableForBooking(context.Context, *types.BookRoomParams, primitive.ObjectID) (bool, error)
}

type MongoBookStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookStore(client *mongo.Client) *MongoBookStore {
	return &MongoBookStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(bookingColl),
	}
}

func (s *MongoBookStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping booking collection")
	return s.coll.Drop(ctx)
}

func (s *MongoBookStore) IsRoomAvailableForBooking(ctx context.Context, newBooking *types.BookRoomParams, roomID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"roomID": roomID,
		"$or": []bson.M{
			// Check if the proposed booking's FromDate falls within the period of existing bookings
			{"fromDate": bson.M{"$lte": newBooking.FromDate}, "tillDate": bson.M{"$gt": newBooking.FromDate}},
			// Check if the proposed booking's TillDate falls within the period of existing bookings
			{"fromDate": bson.M{"$lt": newBooking.TillDate}, "tillDate": bson.M{"$gte": newBooking.TillDate}},
			// Check if the proposed booking's period completely encloses an existing booking's period
			{"fromDate": bson.M{"$gte": newBooking.FromDate}, "tillDate": bson.M{"$lte": newBooking.TillDate}},
		},
	}

	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error counting bookings: %v", err)
	}

	return count == 0, nil
}

func (s *MongoBookStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking *types.Booking
	filter := bson.M{"_id": oid}
	if err := s.coll.FindOne(ctx, filter).Decode(&booking); err != nil {
		return nil, err
	}
	return booking, nil
}
