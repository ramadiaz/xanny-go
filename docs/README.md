# API Documentation

This directory contains the auto-generated Swagger/OpenAPI documentation for the Xanny Go Template API.

## Overview

The API documentation is generated using [swaggo/swag](https://github.com/swaggo/swag), which automatically creates OpenAPI 3.0 specifications from Go code comments.

## Files

- `docs.go` - Main documentation file (auto-generated)
- `swagger.json` - OpenAPI specification in JSON format (auto-generated)
- `swagger.yaml` - OpenAPI specification in YAML format (auto-generated)

## Generating Documentation

To generate the Swagger documentation, run:

```bash
make swagger
```

Or manually:

```bash
swag init -g cmd/server/main.go -o docs
```

## Viewing Documentation

Once the server is running, you can access the Swagger UI at:

```
http://localhost:8080/api/swagger/index.html
```

## API Endpoints

### Health Check
- `GET /api/health` - Check the health status of all services

### User Management
- `POST /api/user/create` - Create a new user account
- `POST /api/user/login` - Authenticate user and get tokens
- `POST /api/user/refresh` - Refresh access token
- `POST /api/user/logout` - Logout user and invalidate tokens

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

1. Include the `Authorization` header with the format: `Bearer <access_token>`
2. The access token is obtained from the login endpoint
3. Use the refresh token endpoint to get a new access token when it expires

## Request/Response Examples

### User Registration
```json
POST /api/user/create
{
  "email": "user@example.com",
  "password": "password123"
}
```

### User Login
```json
POST /api/user/login
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Health Check Response
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "services": {
    "database": {
      "status": "healthy",
      "message": "Database connection successful",
      "latency": "1.234ms"
    },
    "redis": {
      "status": "healthy",
      "message": "Redis connection successful",
      "latency": "0.567ms"
    }
  }
}
```

## Error Responses

The API returns consistent error responses in the following format:

```json
{
  "status": 400,
  "message": "Bad request",
  "error": "Validation failed"
}
```

## Development

When adding new endpoints:

1. Add Swagger comments to your controller functions
2. Use the `@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, and `@Router` annotations
3. Regenerate the documentation using `make swagger`
4. Test the documentation in the Swagger UI

## Example Swagger Comments

```go
// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User registration data"
// @Success 201 {object} dto.Response
// @Failure 400 {object} exceptions.Exception
// @Failure 409 {object} exceptions.Exception
// @Router /user/create [post]
func (h *Handler) CreateUser(ctx *gin.Context) {
    // Implementation
}
``` 