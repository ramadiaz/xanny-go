package middleware

import (
	"net/url"
	"xanny-go-template/models"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"gorm.io/gorm"
)

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

		data := models.Client{
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
