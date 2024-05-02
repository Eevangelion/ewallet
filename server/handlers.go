package server

import (
	"net/http"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	c.JSON(http.StatusOK, wal)
}

func (w *WalletServer) SendMoney(c *gin.Context) {
	logger := logger.GetLogger()
	walId := c.Param("walletId")
	var req contracts.RequestSendMoney
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(
			"Error while parsing body:",
			zap.String("event", "parse_body"),
			zap.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := w.service.BalanceTransfer(walId, req.To, req.Amount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func (w *WalletServer) GetWalletHistory(c *gin.Context) {
	walId := c.Param("walletId")
	hist, err := w.service.GetHistory(walId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	response := contracts.TransactionHistory{}

	response.TransactionList = hist.TransactionList

	c.JSON(http.StatusOK, response.TransactionList)
}

func (w *WalletServer) GetWalletState(c *gin.Context) {
	walId := c.Param("walletId")
	w.service.GetState(walId)
	c.Status(200)
}
