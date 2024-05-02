package logger

import (
	"os"

	"go.uber.org/zap"
)

var logger *zap.Logger

func GetLogger() *zap.Logger {
	if logger == nil {
		if os.Getenv("TESTING") == "true" {
			logger = zap.NewNop()
		} else {
			logger, _ = zap.NewProduction()
		}
	}
	defer logger.Sync()
	return logger
}
