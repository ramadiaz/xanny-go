# Use the official Golang image as a base
FROM golang:1.22.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o main cmd/xanny-go-template/main.go

# Expose the port on which the app runs
EXPOSE 8013

# Command to run the executable
CMD ["./main"]
