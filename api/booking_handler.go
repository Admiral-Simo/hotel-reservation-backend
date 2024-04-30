package api

import (
	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	var (
		page  int64 = parsePageQueryParam(c.Query("page"))
		limit int64 = 10
	)
	opts := options.FindOptions{}
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	bookings, err := h.store.Booking.GetBookings(c.Context(), nil, &opts)
	if err != nil {
		return ErrNotFound("users")
	}
	if bookings == nil {
        return c.JSON([]struct{}{})
	}
    return c.JSON(bookings)
}

// this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAuthorized()
	}
	if err := checkBookingAuth(c, booking, user); err != nil {
		return ErrUnAuthorized()
	}
	// TODO: update booking.Canceled = true
	return c.JSON(booking)
}

// this needs to be user authorized
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAuthorized()
	}
	if err := checkBookingAuth(c, booking, user); err != nil {
		return ErrUnAuthorized()
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
