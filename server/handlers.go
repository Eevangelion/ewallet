package server

import (
	"net/http"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/errs"
	"github.com/Eevangelion/ewallet/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	DefaultBalance float32 = 100
)

func (w *WalletServer) CreateWallet(c *gin.Context) {
	logger := logger.GetLogger()
	wal, err := w.service.Create(DefaultBalance)
	if err != nil {
		err = errs.WrapErr(err, "CreateWallet:")
		logger.Error(
			"Error while creating wallet:",
			zap.String("event", err.Event),
			zap.String("error", err.Message),
		)
		c.JSON(err.Code, nil)
		return
	}
	c.JSON(http.StatusOK, wal)
}

func (w *WalletServer) SendMoney(c *gin.Context) {
	logger := logger.GetLogger()
	walId := c.Param("walletId")
	var req contracts.RequestSendMoney
	if e := c.ShouldBindJSON(&req); e != nil {
		err := errs.NewErr(e, "parse_body", "bad data")
		err = errs.WrapErr(err, "SendMoney:")
		logger.Error(
			"Error while parsing body:",
			zap.String("event", err.Event),
			zap.String("error", err.Message),
		)
		c.JSON(err.Code, gin.H{"error": err.Event})
		return
	}

	err := w.service.BalanceTransfer(walId, req.To, req.Amount)

	if err != nil {
		err = errs.WrapErr(err, "SendMoney:")
		logger.Error(
			"Error while send money:",
			zap.String("event", err.Event),
			zap.String("error", err.Message),
		)
		c.JSON(err.Code, gin.H{"error": err.Event})
		return
	}

	c.Status(200)
}

func (w *WalletServer) GetWalletHistory(c *gin.Context) {
	logger := logger.GetLogger()
	walId := c.Param("walletId")
	hist, err := w.service.GetHistory(walId)

	if err != nil {
		err = errs.WrapErr(err, "GetWalletHistory:")
		logger.Error(
			"Error while getting wallet history:",
			zap.String("event", err.Event),
			zap.String("error", err.Message),
		)
		c.JSON(err.Code, gin.H{"error": err.Event})
		return
	}

	response := contracts.TransactionHistory{}

	response.TransactionList = hist.TransactionList

	c.JSON(http.StatusOK, response.TransactionList)
}

func (w *WalletServer) GetWalletState(c *gin.Context) {
	logger := logger.GetLogger()
	walId := c.Param("walletId")
	wallet, err := w.service.GetState(walId)

	if err != nil {
		err = errs.WrapErr(err, "GetWalletState:")
		logger.Error(
			"Error while getting wallet state:",
			zap.String("event", err.Event),
			zap.String("error", err.Message),
		)
		c.JSON(err.Code, gin.H{"error": err.Event})
		return
	}

	c.JSON(http.StatusOK, wallet)
}
