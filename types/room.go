package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Size    string             `bson:"size" json:"size"`
	SeaSide bool               `bson:"seaSide" json:"seaSide"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
