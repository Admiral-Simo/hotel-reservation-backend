package api

import (
	"context"
	"testing"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testDBName = "hotel-reservation-test"

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {

	err := tdb.client.Database(db.DBNAME).Drop(context.TODO())
	assert.NoError(t, err, "error occurred while dropping database: %v", err)
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	assert.NoError(t, err, "unexpected error connecting to mongodb:%v", err)

	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookStore(client),
		},
	}
}
