package middleware

import (
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"strings"
	"layered-template/dto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error getting secret"})
			return
		}

		var secretKey = []byte(secret)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}

		tokenString := authHeaderParts[1]
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}

		user := dto.User{
			ID:        int64(claims["id"].(float64)),
			Email:     claims["email"].(string),
			Username:  claims["username"].(string),
			FirstName: claims["first_name"].(string),
			LastName:  claims["last_name"].(string),
			Contact:   claims["contact"].(string),
			Address:   claims["address"].(string),
		}

		c.Set("user", user)

		c.Next()
	}
}

func ClientTracker(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		userAgent := c.Request.Header.Get("User-Agent")
		ua := user_agent.New(userAgent)
		name, version := ua.Browser()

		referer := c.Request.Referer()

		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		fullURL := url.URL{
			Path:     path,
			RawQuery: rawQuery,
		}

		_, err := db.Exec("INSERT INTO client_track (ip, browser, version, os, device, origin, api) VALUES($1, $2, $3, $4, $5, $6, $7)", clientIP, name, version, ua.OS(), ua.Platform(), referer, fullURL.String())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}

func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Next()
	}
}
