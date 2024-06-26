package api

import (
	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

// rooms of a specific hotel
func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidId()
	}

	filter := bson.M{"hotelID": oid}
	var (
		page  int64 = parsePageQueryParam(c.Query("page"))
		limit int64 = 10
	)
	opts := options.FindOptions{}
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	rooms, err := h.store.Room.GetRooms(c.Context(), filter, &opts)
	if err != nil {
		return ErrNotFound("rooms")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidId()
	}
	hotel, err := h.store.Hotel.GetHotelById(c.Context(), oid)
	if err != nil {
		return ErrNotFound("hotel")
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var (
		page  int64 = parsePageQueryParam(c.Query("page"))
		limit int64 = 10
	)
	opts := options.FindOptions{}
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil, &opts)
	if err != nil {
		return ErrNotFound("hotels")
	}
	if hotels == nil {
		return c.JSON([]struct{}{})
	}
	return c.JSON(hotels)
}
