package api

import (
	"net/http"
	"strconv"

	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "unauthorized",
		})
	}
	return user, nil
}

func checkBookingAuth(c *fiber.Ctx, booking *types.Booking, user *types.User) error {
	if !user.IsAdmin && booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "unauthorized",
		})
	}
	return nil
}

func parsePageQueryParam(pageString string) int64 {
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		page = 1
	}
	return page
}
