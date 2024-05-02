package server

import (
	"github.com/Eevangelion/ewallet/db/repositories"
	"github.com/Eevangelion/ewallet/service"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine = nil
var s *WalletServer = nil

func GetRouter() *gin.Engine {
	if r != nil {
		return r
	}
	r = gin.Default()

	repo := &repositories.WalletRepository{}
	service := service.WalletService{Repo: repo}
	s = NewWalletServer(&service)

	r.POST("/api/v1/wallet", s.CreateWallet)
	r.POST("api/v1/wallet/:walletId/send", s.SendMoney)
	r.GET("/api/v1/wallet/:walletId/history", s.GetWalletHistory)
	r.GET("/api/v1/wallet/:walletId", s.GetWalletState)
	return r
}
