package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mini-api-golang/config"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/handler"
	"mini-api-golang/internal/middleware"
	"mini-api-golang/internal/routes"
	"mini-api-golang/internal/service"
	"mini-api-golang/pkg/logger"
)

func main() {
	// Initialize logger
	logger.InitializeLogger()
	log := logger.Logger
	defer log.Sync() //nolint:errcheck // Sync errors are non-critical during shutdown

	// Load configuration
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}
	log.Info("configuration loaded", zap.Int("port", cfg.Server.Port))

	// Ensure upload directory exists
	if err := os.MkdirAll(cfg.Upload.Dir, 0750); err != nil {
		log.Fatal("failed to create upload directory", zap.Error(err))
	}

	// Initialize database
	db, err := dao.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatal("failed to initialize database", zap.Error(err))
	}
	log.Info("database initialized", zap.String("path", cfg.Database.Path))

	// Wire up DAOs
	userDAO := dao.NewUserDAO(db)
	taskDAO := dao.NewTaskDAO(db)
	postDAO := dao.NewPostDAO(db)
	likeDAO := dao.NewLikeDAO(db)
	repostDAO := dao.NewRepostDAO(db)
	notificationDAO := dao.NewNotificationDAO(db)
	followDAO := dao.NewFollowDAO(db)
	userSettingsDAO := dao.NewUserSettingsDAO(db)
	postImageDAO := dao.NewPostImageDAO(db)
	commentDAO := dao.NewCommentDAO(db)

	// Wire up services
	userSvc := service.NewUserService(userDAO, userSettingsDAO)
	postSvc := service.NewPostService(postDAO, followDAO, notificationDAO, userSettingsDAO)
	avatarSvc := service.NewAvatarService(userDAO, cfg.Upload.Dir, cfg.Upload.MaxSizeMB)
	postImageSvc := service.NewPostImageService(postImageDAO, cfg.Upload.Dir, cfg.Upload.MaxSizeMB)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notificationDAO)
	repostSvc := service.NewRepostService(repostDAO, postDAO, notificationDAO)
	notificationSvc := service.NewNotificationService(notificationDAO)
	followSvc := service.NewFollowService(followDAO, userDAO, notificationDAO, userSettingsDAO)
	commentSvc := service.NewCommentService(commentDAO, postDAO, userSettingsDAO, followDAO)
	userSettingsSvc := service.NewUserSettingsService(userSettingsDAO)

	// Wire up handlers
	userH := handler.NewUserHandler(userSvc, cfg)
	taskH := handler.NewTaskHandler(taskDAO)
	postH := handler.NewPostHandler(postSvc, postImageSvc, likeSvc, repostSvc)
	avatarH := handler.NewAvatarHandler(avatarSvc, cfg)
	likeH := handler.NewLikeHandler(likeSvc)
	repostH := handler.NewRepostHandler(repostSvc)
	notificationH := handler.NewNotificationHandler(notificationSvc)
	followH := handler.NewFollowHandler(followSvc)
	commentH := handler.NewCommentHandler(commentSvc)
	userSettingsH := handler.NewUserSettingsHandler(userSettingsSvc)

	// Set up Gin engine
	if cfg.Log.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(log))

	routes.SetupRoutes(r, userH, taskH, postH, avatarH, likeH, repostH, notificationH, followH, commentH, userSettingsH, cfg.Upload.Dir, cfg.JWT.Secret)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// Start server in a goroutine so it doesn't block graceful shutdown handling
	go func() {
		log.Info("starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited")
}
