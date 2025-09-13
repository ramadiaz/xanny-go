package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"time"
	"xanny-go/models"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"gorm.io/gorm"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func colorMethod(method string) string {
	switch method {
	case "GET":
		return Blue + method + Reset
	case "POST":
		return Green + method + Reset
	case "PUT":
		return Yellow + method + Reset
	case "DELETE":
		return Red + method + Reset
	default:
		return Cyan + method + Reset
	}
}

func colorStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return Green + fmt.Sprintf("%d", status) + Reset
	case status >= 400 && status < 500:
		return Yellow + fmt.Sprintf("%d", status) + Reset
	case status >= 500:
		return Red + fmt.Sprintf("%d", status) + Reset
	default:
		return fmt.Sprintf("%d", status)
	}
}

func ClientTracker(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
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

			data := models.Clients{
				IP:      clientIP,
				Browser: name,
				Version: version,
				OS:      ua.OS(),
				Device:  ua.Platform(),
				Origin:  referer,
				API:     fullURL.String(),
			}

			db.Create(&data)
		}()
	}
}

func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var reqBody []byte
		if c.Request.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(io.LimitReader(c.Request.Body, 1024))
			reqBody = bodyBytes
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[REQ] %s %s | Body: %s", colorMethod(c.Request.Method), c.Request.URL.Path, truncateForLog(reqBody))
		log.Printf("[RES] %s %s | Status: %s | Duration: %v | Body: %s",
			colorMethod(c.Request.Method), c.Request.URL.Path, colorStatus(status), duration, truncateForLog(blw.body.Bytes()))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func truncateForLog(b []byte) string {
	if len(b) > 1024 {
		return string(b[:1024]) + "..."
	}
	return string(b)
}
