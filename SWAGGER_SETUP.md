# Swagger/OpenAPI Documentation Setup

This document describes the complete Swagger/OpenAPI documentation setup for the Xanny Go Template API.

## What Was Implemented

### 1. Dependencies Added
- `github.com/swaggo/swag` - Swagger documentation generator
- `github.com/swaggo/gin-swagger` - Gin middleware for serving Swagger UI
- `github.com/swaggo/files` - Static files for Swagger UI

### 2. Files Created/Modified

#### New Files:
- `docs/docs.go` - Auto-generated Swagger documentation
- `docs/swagger.json` - OpenAPI specification in JSON format
- `docs/swagger.yaml` - OpenAPI specification in YAML format
- `docs/README.md` - Documentation guide
- `SWAGGER_SETUP.md` - This setup guide

#### Modified Files:
- `cmd/server/main.go` - Added Swagger metadata and docs import
- `routers/main_router.go` - Added Swagger UI endpoint
- `api/users/controllers/users_ctrl_impl.go` - Added Swagger comments to all endpoints
- `api/users/dto/request.go` - Added documentation comments and examples
- `api/users/dto/response.go` - Added documentation comments and examples
- `pkg/helpers/health_helper.go` - Added documentation comments and examples
- `makefile` - Added Swagger generation commands

### 3. API Endpoints Documented

#### Health Check
- `GET /api/health` - System health monitoring

#### User Management
- `POST /api/user/create` - User registration
- `POST /api/user/login` - User authentication
- `POST /api/user/refresh` - Token refresh
- `POST /api/user/logout` - User logout

### 4. Swagger Comments Added

#### Main Application
```go
// @title Xanny Go Template API
// @version 1.0
// @description A comprehensive Go API template with authentication, user management, and health monitoring.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
```

#### Endpoint Documentation
```go
// Create godoc
// @Summary Create a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.Users true "User registration data"
// @Success 201 {object} dto.Response
// @Failure 400 {object} exceptions.Exception
// @Router /user/create [post]
```

#### DTO Documentation
```go
// Users represents user registration request
type Users struct {
    Email     string `json:"email" example:"user@example.com" binding:"required,email"`
    Passoword string `json:"password" example:"password123" binding:"required,min=6"`
}
```

### 5. Makefile Commands Added

```makefile
# Generate Swagger documentation
swagger:
    swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Generate Swagger documentation and run server
swagger-run: swagger
    go run cmd/server/main.go
```

## How to Use

### 1. Generate Documentation
```bash
make swagger
```

### 2. View Documentation
Start the server and visit:
```
http://localhost:8080/api/swagger/index.html
```

### 3. Development Workflow
1. Add Swagger comments to new endpoints
2. Run `make swagger` to regenerate documentation
3. Test in Swagger UI
4. Commit changes

## Features Included

### 1. Interactive API Documentation
- Full OpenAPI 3.0 specification
- Interactive Swagger UI
- Request/response examples
- Authentication documentation

### 2. Comprehensive Coverage
- All user management endpoints
- Health check endpoint
- Error responses
- Request/response schemas

### 3. Developer Experience
- Auto-generated documentation
- Easy regeneration with make commands
- Clear examples and descriptions
- Proper tagging and categorization

### 4. Authentication Documentation
- JWT Bearer token authentication
- Security definitions
- Token refresh flow documentation

## File Structure

```
docs/
├── docs.go          # Auto-generated main documentation
├── swagger.json     # OpenAPI JSON specification
├── swagger.yaml     # OpenAPI YAML specification
└── README.md        # Documentation guide

api/users/
├── controllers/
│   └── users_ctrl_impl.go  # Endpoint documentation
└── dto/
    ├── request.go          # Request DTOs with examples
    └── response.go         # Response DTOs with examples

pkg/helpers/
└── health_helper.go        # Health check documentation
```

## Next Steps

### 1. Add More Endpoints
When adding new endpoints, follow this pattern:
```go
// EndpointName godoc
// @Summary Brief description
// @Description Detailed description
// @Tags tag-name
// @Accept json
// @Produce json
// @Param param-name body dto.TypeName true "Parameter description"
// @Success 200 {object} dto.ResponseType
// @Failure 400 {object} exceptions.Exception
// @Router /endpoint/path [method]
func (h *Handler) EndpointName(ctx *gin.Context) {
    // Implementation
}
```

### 2. Enhance Documentation
- Add more detailed descriptions
- Include more examples
- Document error codes
- Add authentication requirements

### 3. Testing
- Test all endpoints through Swagger UI
- Verify request/response schemas
- Test authentication flows
- Validate error responses

## Troubleshooting

### Common Issues

1. **Swagger not generating**: Run `go mod tidy` first
2. **Missing types**: Ensure all referenced types are properly documented
3. **Import errors**: Use `--parseDependency --parseInternal` flags
4. **UI not loading**: Check if server is running on correct port

### Regeneration
If documentation gets out of sync:
```bash
make swagger
```

## Conclusion

The Swagger/OpenAPI documentation is now fully integrated into the Xanny Go Template API. Developers can:

- View interactive API documentation
- Test endpoints directly from the UI
- Understand request/response formats
- See authentication requirements
- Generate client SDKs from the OpenAPI specification

The documentation will automatically stay in sync with code changes when the `make swagger` command is run. 