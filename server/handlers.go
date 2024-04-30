package server

import (
	"net/http"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/gin-gonic/gin"
)

const (
	DefaultBalance float32 = 100
)

func (w *WalletServer) CreateWallet(c *gin.Context) {
	wal, err := w.service.Create(DefaultBalance)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":      wal.Id,
		"balance": wal.Balance,
	})
}

func (w *WalletServer) SendMoney(c *gin.Context) {
	wal_id := c.Param("walletId")
	var req contracts.RequestSendMoney
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w.service.BalanceTransfer(wal_id, req.To, req.Amount)

	c.Status(200)
}

func (w *WalletServer) GetWalletHistory(c *gin.Context) {
	wal_id := c.Param("walletId")
	w.service.GetWalletHistory(wal_id)
	c.Status(200)
}

func (w *WalletServer) GetWalletState(c *gin.Context) {
	wal_id := c.Param("walletId")
	w.service.GetWalletState(wal_id)
	c.Status(200)
}
