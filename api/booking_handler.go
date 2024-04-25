package api

import (
	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// this needs to be admin authorized!
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	if bookings == nil {
		return c.JSON([]types.Booking{})
	}
	return c.JSON(bookings)
}

// this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if err := checkBookingAuth(c, booking, user); err != nil {
		return err
	}
	// TODO: update booking.Canceled = true
	return c.JSON(booking)
}

// this needs to be user authorized
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if err := checkBookingAuth(c, booking, user); err != nil {
		return err
	}
    update := bson.M{"$set": bson.M{"canceled": true}}
	if err := h.store.Booking.UpdateBookingById(c.Context(), id, update); err != nil {
		return err
	}
	return c.JSON(genericResponse{
		Type: "success",
		Msg:  "booking canceled successfuly",
	})
}
