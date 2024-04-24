package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID     primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons int                `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	FromDate   time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate   time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p *BookRoomParams) CreateBooking(userID, roomID primitive.ObjectID) *Booking {
	return &Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: p.NumPersons,
		FromDate:   p.FromDate,
		TillDate:   p.TillDate,
	}
}

func (p *BookRoomParams) Validate() error {
	now := time.Now()
	if p.FromDate.Before(now) {
		return fmt.Errorf("cannot book a room in the past")
	}
	duration := p.TillDate.Sub(p.FromDate)
	if duration < 24*time.Hour {
		return fmt.Errorf("invalid booking duration: should be at least one day")
	}
	return nil
}
