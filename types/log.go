package types

// `from` `route` `Header` `method` `body`

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	From      string                 `bson:"from" json:"from"`
	Route     string                 `bson:"route" json:"route"`
	Header    map[string]interface{} `bson:"header" json:"header"`
	Method    string                 `bson:"method" json:"method"`
	Body      map[string]interface{} `bson:"body" json:"body"`
	TimeStamp time.Time              `bson:"timestamp" json:"timestamp"`
}
