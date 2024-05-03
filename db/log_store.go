package db

import (
	"context"

	"github.com/Admiral-Simo/HotelReserver/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	logsColl = "logs"
)

type LogsStore interface {
	Dropper

	InsertLog(context.Context, *types.Log) error
	GetLogs(context.Context, bson.M, *options.FindOptions) ([]*types.Log, error)
	//GetBookingById(ctx context.Context, id string) (*types.Booking, error)
	//IsRoomAvailableForBooking(context.Context, *types.BookRoomParams, primitive.ObjectID) (bool, error)
	//UpdateBookingById(ctx context.Context, id string, update bson.M) error
}

type MongoLogsStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	LogsStore
}

func NewMongoLogsStore(client *mongo.Client) *MongoLogsStore {
	return &MongoLogsStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(logsColl),
	}
}

func (s *MongoLogsStore) InsertLog(ctx context.Context, log *types.Log) error {
	_, err := s.coll.InsertOne(ctx, log)
	return err
}

func (s *MongoLogsStore) GetLogs(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*types.Log, error) {
	cur, err := s.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var logs []*types.Log
	if err = cur.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
