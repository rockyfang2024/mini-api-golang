package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"mini-api-golang/internal/handler"
	"mini-api-golang/internal/middleware"
)

// SetupRoutes registers all application routes on the given gin.Engine.
func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, taskHandler *handler.TaskHandler, postHandler *handler.PostHandler, jwtSecret string) {
	// CORS — allow the Vite dev server and same-origin requests
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Legacy public auth routes (kept for backward compatibility)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Legacy protected user routes
	users := r.Group("/users", middleware.JWTMiddleware(jwtSecret))
	{
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// Legacy protected task routes
	tasks := r.Group("/tasks", middleware.JWTMiddleware(jwtSecret))
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}

	// ── /api routes ──────────────────────────────────────────────────────────

	api := r.Group("/api")

	// Auth (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// Current user info (protected)
	api.GET("/me", middleware.JWTMiddleware(jwtSecret), userHandler.Me)

	// Posts — list is semi-public (uses optional auth to filter visibility)
	// We register two versions of GET /api/posts: one with auth middleware that
	// stores user_id, and one fallback for unauthenticated callers.
	// Instead, we use a single route and an optional-auth helper.
	api.GET("/posts", middleware.OptionalJWT(jwtSecret), postHandler.ListPosts)
	api.POST("/posts", middleware.JWTMiddleware(jwtSecret), postHandler.CreatePost)

	// User posts (optional auth for visibility filtering)
	api.GET("/users/:id/posts", middleware.OptionalJWT(jwtSecret), postHandler.ListUserPosts)
}

