package middleware

import (
	"xanny-go/pkg/exceptions"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limiters ...*limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, lmt := range limiters {
			httpError := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
			if httpError != nil {
				ctx.AbortWithStatusJSON(httpError.StatusCode, exceptions.NewException(httpError.StatusCode, httpError.Message))
				return
			}
		}
		ctx.Next()
	}
}
