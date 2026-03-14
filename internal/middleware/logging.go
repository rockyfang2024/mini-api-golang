package middleware

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// Logger is a middleware for logging HTTP requests.
func Logger() gin.HandlerFunc {
    logger, _ := zap.NewProduction()
    defer logger.Sync() // flushes buffer, if any

    return func(c *gin.Context) {
        start := time.Now()

        // Process the request
        c.Next()

        // Log the request
        logger.Info("Request details",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", time.Since(start)),
        )
    }
}