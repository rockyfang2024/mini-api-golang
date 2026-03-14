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
func SetupRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	taskHandler *handler.TaskHandler,
	postHandler *handler.PostHandler,
	avatarHandler *handler.AvatarHandler,
	likeHandler *handler.LikeHandler,
	repostHandler *handler.RepostHandler,
	notificationHandler *handler.NotificationHandler,
	followHandler *handler.FollowHandler,
	commentHandler *handler.CommentHandler,
	userSettingsHandler *handler.UserSettingsHandler,
	uploadDir string,
	jwtSecret string,
) {
	// CORS — allow the Vite dev server and same-origin requests
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve uploaded files as static assets under /uploads
	r.Static("/uploads", uploadDir)

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

	// Current user info + avatar upload (protected)
	api.GET("/me", middleware.JWTMiddleware(jwtSecret), userHandler.Me)
	api.POST("/me/avatar", middleware.JWTMiddleware(jwtSecret), avatarHandler.UploadAvatar)
	api.GET("/settings", middleware.JWTMiddleware(jwtSecret), userSettingsHandler.GetSettings)
	api.PUT("/settings", middleware.JWTMiddleware(jwtSecret), userSettingsHandler.UpdateSettings)

	// Posts — list is semi-public (uses optional auth to filter visibility)
	api.GET("/posts", middleware.OptionalJWT(jwtSecret), postHandler.ListPosts)
	api.POST("/posts", middleware.JWTMiddleware(jwtSecret), postHandler.CreatePost)
	api.GET("/posts/:id/comments", middleware.OptionalJWT(jwtSecret), commentHandler.ListComments)
	api.POST("/posts/:id/comments", middleware.JWTMiddleware(jwtSecret), commentHandler.CreateComment)

	// Post interactions (protected)
	postGroup := api.Group("/posts", middleware.JWTMiddleware(jwtSecret))
	{
		postGroup.POST("/:id/like", likeHandler.LikePost)
		postGroup.DELETE("/:id/like", likeHandler.UnlikePost)
		postGroup.POST("/:id/repost", repostHandler.RepostPost)
	}

	api.POST("/comments/:id/replies", middleware.JWTMiddleware(jwtSecret), commentHandler.ReplyComment)

	// User posts (optional auth for visibility filtering)
	api.GET("/users/:id/posts", middleware.OptionalJWT(jwtSecret), postHandler.ListUserPosts)
	api.GET("/users/:id", middleware.OptionalJWT(jwtSecret), userHandler.GetUserProfile)

	// Follow routes (protected for mutation, public for reading)
	api.POST("/users/:id/follow", middleware.JWTMiddleware(jwtSecret), followHandler.FollowUser)
	api.DELETE("/users/:id/follow", middleware.JWTMiddleware(jwtSecret), followHandler.UnfollowUser)
	api.GET("/users/:id/followers", followHandler.ListFollowers)
	api.GET("/users/:id/following", followHandler.ListFollowing)

	// Notifications (all protected)
	notifs := api.Group("/notifications", middleware.JWTMiddleware(jwtSecret))
	{
		notifs.GET("", notificationHandler.ListNotifications)
		notifs.PUT("/read-all", notificationHandler.MarkAllNotificationsRead)
		notifs.PUT("/:id/read", notificationHandler.MarkNotificationRead)
	}
}
