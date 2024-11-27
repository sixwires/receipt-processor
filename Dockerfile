# Use the official Golang image as the base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the rest of the application files
COPY . .

# Build the app
RUN go build -v -o receipt-processor ./...

# Run the executable
CMD ["./receipt-processor"]
