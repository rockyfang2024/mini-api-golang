package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mini-api-golang/internal/handler"
	"mini-api-golang/internal/middleware"
)

// SetupRoutes registers all application routes on the given gin.Engine.
func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, taskHandler *handler.TaskHandler, jwtSecret string) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public auth routes
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Protected user routes
	users := r.Group("/users", middleware.JWTMiddleware(jwtSecret))
	{
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// Protected task routes
	tasks := r.Group("/tasks", middleware.JWTMiddleware(jwtSecret))
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}
}
