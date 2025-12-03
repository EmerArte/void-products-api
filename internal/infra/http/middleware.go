package http

import (
	"strings"
	"time"

	"github.com/emerarteaga/products-api/internal/config"
	"github.com/emerarteaga/products-api/internal/infra/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		logger.Info("HTTP request", "method", method, "path", path, "status", statusCode, "latency", latency.String(), "ip", c.ClientIP())
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("panic recovered", "error", recovered)
		c.JSON(500, gin.H{"success": false, "error": "Internal server error"})
	})
}

// CORS returns a middleware that handles CORS with configuration
func CORS(corsConfig config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowedOrigin := "*"
		if len(corsConfig.AllowedOrigins) > 0 {
			// If wildcard is in the list, allow all origins
			hasWildcard := false
			for _, allowedOrigin := range corsConfig.AllowedOrigins {
				if allowedOrigin == "*" {
					hasWildcard = true
					break
				}
			}

			if hasWildcard {
				allowedOrigin = "*"
			} else if origin != "" {
				// Check if the origin is in the allowed list
				for _, allowed := range corsConfig.AllowedOrigins {
					if allowed == origin {
						allowedOrigin = origin
						break
					}
				}
				if allowedOrigin == "*" && origin != "" {
					// Origin not in allowed list, deny
					c.AbortWithStatus(403)
					return
				}
			}
		}

		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
