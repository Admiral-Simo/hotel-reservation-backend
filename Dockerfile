# Use the official Golang image with Alpine Linux as a base
FROM golang:1.22.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 3000 for the Go application
EXPOSE 3000

# Start the MongoDB service and then run the Go application
CMD ./main
