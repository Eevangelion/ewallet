package main

import (
	"fmt"

	"github.com/Eevangelion/ewallet/config"
	"github.com/Eevangelion/ewallet/db"
	"github.com/Eevangelion/ewallet/logger"
	"github.com/Eevangelion/ewallet/server"
	"go.uber.org/zap"
)

func main() {
	conf := config.GetConfig()
	r := server.GetRouter()
	logger := logger.GetLogger()
	pool, err := db.GetPool()
	if err != nil {
		logger.Error(
			"Error while connecting to DB:",
			zap.String("event", "connect_database"),
			zap.String("error", err.Error()),
		)
	}
	defer pool.Close()
	port := conf.Server.Port
	r.Run(fmt.Sprintf(":%d", port))
}