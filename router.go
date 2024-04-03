package main

import (
	"github.com/Eevangelion/ewallet/handlers"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine = nil

func GetRouter() *gin.Engine {
	if r != nil {
		return r
	}
	r = gin.Default()

	r.POST("/api/v1/wallet", handlers.CreateWallet)
	r.POST("api/v1/wallet/:walletId/send", handlers.SendMoney)
	r.GET("/api/v1/wallet/:walletId/history", handlers.GetWalletHistory)
	r.GET("/api/v1/wallet/:walletId", handlers.GetWalletState)
	return r
}
