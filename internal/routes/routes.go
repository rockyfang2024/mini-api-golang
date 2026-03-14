package routes

import (
    "github.com/gin-gonic/gin"
)

// SetupRoutes initializes the routes for the API
func SetupRoutes(r *gin.Engine) {
    r.GET("/health", HealthCheck)
    // Add more routes here
}

// HealthCheck is a simple health check endpoint
func HealthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
    })
}

// Middleware example
func ExampleMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Log or manipulate request here
        c.Next()
        // Log or manipulate response here
    }
}
