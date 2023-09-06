# Use the official Golang image from the Docker Hub
FROM golang:1.21

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

COPY .env .env

# Build the Go app
RUN go build -o main .

# Expose port 8080 to interact with the application
EXPOSE 8080

# Command to run the application
CMD ["./main"]