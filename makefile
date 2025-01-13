# Run the application
run:
	go mod tidy
	go run cmd/xanny-go-template/main.go

# Build the application binary
build:
	go build -o bin/xanny-go-template ./cmd/xanny-go-template

# Run migrations
migrate:
	go run cmd/migrate/migrate.go

# Clean the build (remove binaries and build artifacts)
clean:
	rm -f bin/xanny-go-template