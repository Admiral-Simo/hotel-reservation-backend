package api

import (
	"time"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookRoomParams struct {
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	var (
		page  int64 = parsePageQueryParam(c.Query("page"))
		limit int64 = 10
	)
	opts := options.FindOptions{}
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	rooms, err := h.store.Room.GetRooms(c.Context(), nil, &opts)
	if err != nil {
		return err
	}
	if rooms == nil {
		return c.JSON([]struct{}{})
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleGetRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	rooms, err := h.store.Room.GetRoomById(c.Context(), oid)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params *types.BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	roomID := c.Params("id")
	roomOID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return ErrInvalidId()
	}

	user, err := getAuthUser(c)

	if err != nil {
		return ErrUnAuthorized()
	}

	if err = params.Validate(); err != nil {
		return ErrBadRequest()
	}

	ok, err := h.store.Booking.IsRoomAvailableForBooking(c.Context(), params, roomOID)
	if err != nil {
		return ErrNotFound("room")
	}

	if !ok {
		return ErrUnavailable("room")
	}

	booking := params.CreateBooking(user.ID, roomOID)

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), booking)

	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}
