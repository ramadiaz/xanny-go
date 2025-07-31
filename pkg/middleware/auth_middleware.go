package middleware

import (
	"net/http"
	"os"
	"strings"
	"xanny-go-template/api/users/dto"
	"xanny-go-template/pkg/exceptions"
	"xanny-go-template/pkg/helpers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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

		isBlacklisted, _ := helpers.IsTokenBlacklisted(tokenString)
		if isBlacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NewException(http.StatusUnauthorized, "Token is blacklisted"))
			return
		}

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

		user := dto.UserOutput{
			UUID:  claims["uuid"].(string),
			Email: claims["email"].(string),
			Name:  claims["name"].(string),
		}

		c.Set("user", user)
		c.Next()
	}
}
