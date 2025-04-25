# -------- STAGE 1: Build --------
FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/xanny-go-template/main.go


# -------- STAGE 2: Runtime --------
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy only the final binary from the builder stage
COPY --from=builder /app/main .

# Expose the app port (set at runtime via env var)
EXPOSE ${PORT}

# Run the binary
CMD ["./main"]
