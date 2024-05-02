package logger

import "go.uber.org/zap"

var logger *zap.Logger

func GetLogger() *zap.Logger {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	return logger
}
