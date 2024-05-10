package main

import (
	"context"
	"log"
	"os"

	"github.com/Admiral-Simo/HotelReserver/api"
	"github.com/Admiral-Simo/HotelReserver/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	// load envirement variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("couldn't load envirement variables")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// stores
	// handlers
	// api applications
	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookStore(client)
		logsStore    = db.NewMongoLogsStore(client)
		store        = &db.Store{
			User:    userStore,
			Room:    roomStore,
			Hotel:   hotelStore,
			Booking: bookingStore,
		}
	)

	app := fiber.New(config)

	// Middleware for CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	app.Use(api.Logger(logsStore))

	// groups
	var (
		apiv1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin = apiv1.Group("/admin", api.AdminAuth)
		auth  = app.Group("/api")
	)

	// handlers
	var (
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versions api routes
	// user routes
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel routes
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Get("/room/:id", roomHandler.HandleGetRoom)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	// TODO: Cancel a booking

	// bookings routes
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin routes
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	admin.Get("/user", userHandler.HandleGetUsers)

	app.Use(api.NotFoundHandler)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")

	app.Listen(listenAddr)
}
