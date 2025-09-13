package helpers

import (
	"context"
	"xanny-go/pkg/config"
	"time"

	"gorm.io/gorm"
)

type HealthCheck struct {
	Status    string                   `json:"status" example:"healthy"`
	Timestamp time.Time                `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Services  map[string]ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Status  string `json:"status" example:"healthy"`
	Message string `json:"message,omitempty" example:"Database connection successful"`
	Latency string `json:"latency,omitempty" example:"1.234ms"`
}

func CheckDatabaseHealth(db *gorm.DB) ServiceStatus {
	start := time.Now()

	sqlDB, err := db.DB()
	if err != nil {
		return ServiceStatus{
			Status:  "error",
			Message: "Failed to get database instance: " + err.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	latency := time.Since(start).String()

	if err != nil {
		return ServiceStatus{
			Status:  "error",
			Message: "Database connection failed: " + err.Error(),
			Latency: latency,
		}
	}

	return ServiceStatus{
		Status:  "healthy",
		Message: "Database connection successful",
		Latency: latency,
	}
}

func CheckRedisHealth() ServiceStatus {
	start := time.Now()

	if config.RedisClient == nil {
		return ServiceStatus{
			Status:  "error",
			Message: "Redis client not initialized",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.RedisClient.Ping(ctx).Result()
	latency := time.Since(start).String()

	if err != nil {
		return ServiceStatus{
			Status:  "error",
			Message: "Redis connection failed: " + err.Error(),
			Latency: latency,
		}
	}

	return ServiceStatus{
		Status:  "healthy",
		Message: "Redis connection successful",
		Latency: latency,
	}
}

func PerformHealthCheck(db *gorm.DB) HealthCheck {
	services := make(map[string]ServiceStatus)

	services["database"] = CheckDatabaseHealth(db)
	services["redis"] = CheckRedisHealth()

	overallStatus := "healthy"
	for _, service := range services {
		if service.Status == "error" {
			overallStatus = "unhealthy"
			break
		}
	}

	return HealthCheck{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
	}
}
