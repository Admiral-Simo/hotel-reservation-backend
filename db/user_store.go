package db

import (
	"context"
	"fmt"

	"github.com/Admiral-Simo/HotelReserver/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userColl = "users"

type UserStore interface {
	Dropper

	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context, bson.M, *options.FindOptions) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping users collection")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.M{"$set": params}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// TODO: maybe it's a good idea to handle if we did not delete any user
	// maybe log it or something??
	filter := bson.M{"_id": oid}
	_, err = s.coll.DeleteOne(ctx, filter)
	return err
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	filter := bson.M{"email": email}
	var user *types.User
	if err := s.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}
