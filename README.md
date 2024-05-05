# Hotel reservation backend

# environment variables
```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=somethingsupersecretthatNOBODYKNOWS
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
```

## Project outline
- users -> book room from an hotel
- admins -> going to check reservations/bookings
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration

## Resources
### Mongodb driver
Documentation
```
http://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber
Documentation
```
https://gofiber.io
```

Installing gofiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Installing mongodb as a Docker container
```
docker run --name my_mongo -d -p 27017:27017 mongo:latest
```
