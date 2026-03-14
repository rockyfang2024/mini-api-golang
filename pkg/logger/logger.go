package logger

import (
	"go.uber.org/zap"
)

var (
	// Global logger variable
	Logger *zap.Logger
)

// InitializeLogger initializes the zap logger.
func InitializeLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}
}

// Debug logs a message at debug level.
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info logs a message at info level.
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Error logs a message at error level.
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}