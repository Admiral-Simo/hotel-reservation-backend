package api

import (
	"encoding/json"
	"strings"

	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/Admiral-Simo/HotelReserver/types"
	"github.com/gofiber/fiber/v2"
)

// this logger will push to the database
// `from` `route` `Header` `method` `body`

func Logger(logsStore db.LogsStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		from := getClientIP(c)

		headers := parseHeaders(c.Request().Header.Header())

		body := make(map[string]interface{})
		err := json.Unmarshal(c.Body(), &body)
		if err != nil {
			body = make(map[string]interface{})
		}

		log := &types.Log{
			From:   from,
			Route:  c.Path(),
			Method: c.Method(),
			Header: headers,
			Body:   body,
		}

		if err := logsStore.InsertLog(c.Context(), log); err != nil {
			return err
		}

		return c.Next()
	}
}

func getClientIP(c *fiber.Ctx) string {
	ip := c.IP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}

func parseHeaders(headerBytes []byte) map[string]interface{} {
	headers := make(map[string]interface{})
	headerStr := string(headerBytes)
	lines := strings.Split(headerStr, "\r\n")
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}
	return headers
}
