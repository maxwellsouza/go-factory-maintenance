package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// LoggerMiddleware cria logs estruturados em JSON com ID de requisição.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()

		c.Set("request_id", requestID)

		c.Next() // processa a request

		latency := time.Since(start)

		entry := log.WithFields(log.Fields{
			"request_id": requestID,
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    latency,
		})

		if len(c.Errors) > 0 {
			entry.WithField("errors", c.Errors.String()).Error("request failed")
		} else {
			entry.Info("request completed")
		}
	}
}
