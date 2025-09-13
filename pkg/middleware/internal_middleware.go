package middleware

import (
	"net/http"
	"strings"
	"xanny-go/pkg/config"
	"xanny-go/pkg/exceptions"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func InternalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		INTERNAL_SECRET := config.GetInternalSecret()
		ADMIN_USERNAME := config.GetAdminUsername()

		var secretKey = []byte(INTERNAL_SECRET)

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

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, exceptions.NewException(http.StatusForbidden, exceptions.ErrInvalidCredentials))
			return
		}

		username := claims["admin_username"].(string)

		if ADMIN_USERNAME != username {
			c.AbortWithStatusJSON(http.StatusForbidden, exceptions.NewException(http.StatusForbidden, exceptions.ErrInvalidCredentials))
			return
		}

		c.Next()
	}
}
