package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	tokens, ok := c.GetReqHeaders()["X-Api-Token"]

	tokenString := tokens[0]

	if !ok {
		return fmt.Errorf("unauthorized")
	}

	claims, err := validateToken(tokenString)

	if err != nil {
		return err
	}

	expires, err := time.Parse(time.RFC3339, claims["expires"].(string))

	if err != nil {
		return fmt.Errorf("unauthorized")
	}

	if time.Now().After(expires) {
		return fmt.Errorf("token expired")
	}

	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		fmt.Println("Never print a secret for debuging", secret)

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
