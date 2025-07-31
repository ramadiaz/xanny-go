package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheMiddlewareOptions holds configuration for cache middleware
type CacheMiddlewareOptions struct {
	Cache          Cache
	DefaultTTL     time.Duration
	KeyGenerator   func(c *gin.Context) string
	SkipCache      func(c *gin.Context) bool
	IncludeHeaders []string
	ExcludeHeaders []string
	IncludeMethods []string
	ExcludeMethods []string
}

// DefaultCacheMiddlewareOptions returns default options for cache middleware
func DefaultCacheMiddlewareOptions(cache Cache) *CacheMiddlewareOptions {
	return &CacheMiddlewareOptions{
		Cache:          cache,
		DefaultTTL:     5 * time.Minute,
		IncludeMethods: []string{"GET"},
		ExcludeMethods: []string{"POST", "PUT", "DELETE", "PATCH"},
		IncludeHeaders: []string{"Accept", "Accept-Language"},
		ExcludeHeaders: []string{"Authorization", "Cookie"},
	}
}

// CacheMiddleware creates a new cache middleware
func CacheMiddleware(opts *CacheMiddlewareOptions) gin.HandlerFunc {
	if opts == nil {
		opts = DefaultCacheMiddlewareOptions(nil)
	}

	if opts.KeyGenerator == nil {
		opts.KeyGenerator = defaultKeyGenerator
	}

	if opts.SkipCache == nil {
		opts.SkipCache = defaultSkipCache
	}

	return func(c *gin.Context) {
		// Skip caching if configured
		if opts.SkipCache(c) {
			c.Next()
			return
		}

		// Generate cache key
		cacheKey := opts.KeyGenerator(c)

		// Try to get from cache
		if cachedResponse, err := opts.Cache.Get(c.Request.Context(), cacheKey); err == nil {
			var response cachedHTTPResponse
			if err := json.Unmarshal(cachedResponse, &response); err == nil {
				// Set headers
				for key, values := range response.Headers {
					for _, value := range values {
						c.Header(key, value)
					}
				}

				// Set status code
				c.Status(response.StatusCode)

				// Write response body
				c.Data(response.StatusCode, response.ContentType, response.Body)

				// Add cache hit header
				c.Header("X-Cache", "HIT")
				return
			}
		}

		// Create a custom response writer to capture the response
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Cache the response if it's successful
		if c.Writer.Status() == http.StatusOK && len(responseWriter.body.Bytes()) > 0 {
			response := cachedHTTPResponse{
				StatusCode:  c.Writer.Status(),
				ContentType: c.Writer.Header().Get("Content-Type"),
				Body:        responseWriter.body.Bytes(),
				Headers:     make(map[string][]string),
			}

			// Copy headers
			for key, values := range c.Writer.Header() {
				if !isExcludedHeader(key, opts.ExcludeHeaders) {
					response.Headers[key] = values
				}
			}

			// Serialize and cache
			if data, err := json.Marshal(response); err == nil {
				opts.Cache.Set(c.Request.Context(), cacheKey, data, opts.DefaultTTL)
			}
		}

		// Add cache miss header
		c.Header("X-Cache", "MISS")
	}
}

// cachedHTTPResponse represents a cached HTTP response
type cachedHTTPResponse struct {
	StatusCode  int                 `json:"status_code"`
	ContentType string              `json:"content_type"`
	Body        []byte              `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

// responseWriter captures the response for caching
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// defaultKeyGenerator generates a cache key based on request details
func defaultKeyGenerator(c *gin.Context) string {
	// Create a hash of the request
	hash := md5.New()

	// Include method and path
	hash.Write([]byte(c.Request.Method + ":" + c.Request.URL.Path))

	// Include query parameters
	if c.Request.URL.RawQuery != "" {
		hash.Write([]byte("?" + c.Request.URL.RawQuery))
	}

	// Include relevant headers
	for _, header := range []string{"Accept", "Accept-Language", "User-Agent"} {
		if value := c.GetHeader(header); value != "" {
			hash.Write([]byte(":" + header + ":" + value))
		}
	}

	// Include user ID if available
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			hash.Write([]byte(":user:" + id))
		}
	}

	return "http:" + hex.EncodeToString(hash.Sum(nil))
}

// defaultSkipCache determines if caching should be skipped
func defaultSkipCache(c *gin.Context) bool {
	// Skip non-GET requests
	if c.Request.Method != "GET" {
		return true
	}

	// Skip requests with no-cache header
	if c.GetHeader("Cache-Control") == "no-cache" {
		return true
	}

	// Skip requests with authorization header
	if c.GetHeader("Authorization") != "" {
		return true
	}

	return false
}

// isExcludedHeader checks if a header should be excluded from caching
func isExcludedHeader(header string, excludedHeaders []string) bool {
	headerLower := strings.ToLower(header)
	for _, excluded := range excludedHeaders {
		if strings.ToLower(excluded) == headerLower {
			return true
		}
	}
	return false
}

// InvalidateCacheMiddleware creates middleware to invalidate cache
func InvalidateCacheMiddleware(cache Cache, patterns []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request first
		c.Next()

		// Invalidate cache after successful modification
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			for _, pattern := range patterns {
				cache.InvalidatePattern(c.Request.Context(), pattern)
			}
		}
	}
}

// CacheControlMiddleware adds cache control headers
func CacheControlMiddleware(maxAge time.Duration, public bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var directives []string

		if public {
			directives = append(directives, "public")
		} else {
			directives = append(directives, "private")
		}

		if maxAge > 0 {
			directives = append(directives, fmt.Sprintf("max-age=%d", int(maxAge.Seconds())))
		}

		c.Header("Cache-Control", strings.Join(directives, ", "))
		c.Next()
	}
}
