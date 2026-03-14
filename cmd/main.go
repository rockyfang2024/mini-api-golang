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

	// Initialize database
	db, err := dao.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatal("failed to initialize database", zap.Error(err))
	}
	log.Info("database initialized", zap.String("path", cfg.Database.Path))

	// Wire up layers
	userDAO := dao.NewUserDAO(db)
	taskDAO := dao.NewTaskDAO(db)

	userSvc := service.NewUserService(userDAO)

	userH := handler.NewUserHandler(userSvc, cfg)
	taskH := handler.NewTaskHandler(taskDAO)

	// Set up Gin engine
	if cfg.Log.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(log))

	routes.SetupRoutes(r, userH, taskH, cfg.JWT.Secret)

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
