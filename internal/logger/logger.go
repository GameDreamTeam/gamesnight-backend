package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	Logger *zap.Logger
}

var l *Logger

func New() {
	// Handle errors
	logger, _ := zap.NewDevelopment()
	l = &Logger{
		Logger: logger,
	}
}

func GetLogger() *Logger {
	return l
}
