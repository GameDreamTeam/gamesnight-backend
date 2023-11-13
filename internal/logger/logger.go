package logger

import (
	"gamesnight/internal/config"

	"go.uber.org/zap"
)

type Logger struct {
	Logger *zap.Logger
}

// Logger Instance
var l *Logger

func New() {

	var logger *zap.Logger
	var err error

	if config.Get().Env == "local" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic("Unable to iniatilize Logger")
	}

	l = &Logger{
		Logger: logger,
	}

	logger.Info("Logger iniatilized successfully :)")
}

func GetLogger() *Logger {
	return l
}
