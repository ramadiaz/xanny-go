// @title Xanny Go Template API
// @version 1.0
// @description A comprehensive Go API template with authentication, user management, and health monitoring.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xanny-go/docs"
	"xanny-go/pkg/config"
	"xanny-go/pkg/logger"
	"xanny-go/pkg/middleware"
	"xanny-go/routers"

	internalRouters "xanny-go/internal/routers"
	"xanny-go/pkg/helpers"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	docs.SwaggerInfo.BasePath = "/api"
	config.InitEnvCheck()
	helpers.InitRedis()

	logger.Startup()
	port := os.Getenv("PORT")
	environment := os.Getenv("ENVIRONMENT")

	r := gin.New()
	r.Use(middleware.RequestResponseLogger())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	db := config.InitDB()
	validate := validator.New(validator.WithRequiredStructEnabled())
	lmt := tollbooth.NewLimiter(5, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})

	r.Use(middleware.ClientTracker(db))
	r.Use(middleware.GzipResponseMiddleware())
	r.Use(middleware.RateLimitMiddleware(lmt))

	internal := r.Group("/internal")
	internalRouters.InternalRouters(internal, db, validate)

	api := r.Group("/api")
	routers.CompRouters(api, db, validate)

	var host string
	switch environment {
	case "development":
		host = "localhost"
	case "production":
		host = "0.0.0.0"
	default:
		panic("ENV ERROR: {ENVIRONMENT} UNKNOWN")
	}

	server := host + ":" + port

	srv := &http.Server{
		Addr:    server,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server started on port :%s", port)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case sig := <-shutdown:
		log.Printf("Start shutdown... Signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			if err := srv.Close(); err != nil {
				log.Fatalf("Could not force close server: %v", err)
			}
		}
	}

	log.Println("Server stopped")
}
