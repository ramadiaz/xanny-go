package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"xanny-go-template/exceptions"
	"xanny-go-template/models/database"
	"xanny-go-template/models/dto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"gorm.io/gorm"
)

func GzipResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gzipWriter := gzip.NewWriter(c.Writer)
		defer gzipWriter.Close()

		wrappedWriter := &gzipResponseWriter{
			ResponseWriter: c.Writer,
			Writer:         gzipWriter,
		}

		c.Writer = wrappedWriter
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Set("Vary", "Accept-Encoding")

		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("JWT_SECRET")

		var secretKey = []byte(secret)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, exceptions.NewException(http.StatusForbidden, exceptions.ErrForbidden))
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusForbidden, exceptions.NewException(http.StatusForbidden, exceptions.ErrInvalidCredentials))
			return
		}

		tokenString := authHeaderParts[1]
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, exceptions.NewException(http.StatusForbidden, exceptions.ErrInvalidCredentials))
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, exceptions.NewException(http.StatusForbidden, exceptions.ErrInvalidCredentials))
			return
		}

		user := dto.User{
			ID:    claims["id"].(string),
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
		}

		c.Set("user", user)

		c.Next()
	}
}

func ClientTracker(db *gorm.DB) gin.HandlerFunc {
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

		data := database.Client{
			IP:      clientIP,
			Browser: name,
			Version: version,
			OS:      ua.OS(),
			Device:  ua.Platform(),
			Origin:  referer,
			API:     fullURL.String(),
		}

		go db.Create(&data)
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
