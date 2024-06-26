package types

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 8
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *UpdateUserParams) ToBSON() (bson.M, error) {
	m := bson.M{}
	if len(u.FirstName) > 0 {
		if len(u.FirstName) < minFirstNameLen {
			return nil, fmt.Errorf("firstName length should be at least %d characters", minFirstNameLen)
		}
		m["firstName"] = u.FirstName
	}
	if len(u.LastName) > 0 {
		if len(u.LastName) < minLastNameLen {
			return nil, fmt.Errorf("lastName length should be at least %d characters", minFirstNameLen)
		}
		m["lastName"] = u.LastName
	}
	return m, nil
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool
}

func (params *CreateUserParams) Valide() map[string]string {
	errs := make(map[string]string)
	if len(params.FirstName) < minFirstNameLen {
		errs["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errs["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errs["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errs["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func CreateTokenFromUser(user *User) string {
	expires := time.Now().Add(time.Hour * 24)

	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenString
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
		IsAdmin:           params.IsAdmin,
	}, nil
}
