package handlers

import (
	"net/http"
	"strconv"

	"github.com/Eevangelion/ewallet/models"
	"github.com/gin-gonic/gin"
)

type Request struct {
	To     int     `json:"to"`
	Amount float32 `json:"amount"`
}

func CreateWallet(c *gin.Context) {
	wallet := models.Wallet{
		Id:      1,
		Balance: 100,
	}
	func(wal *models.Wallet) {
		// create wallet
	}(&wallet)
	c.JSON(http.StatusOK, gin.H{
		"id":      wallet.Id,
		"balance": wallet.Balance,
	})
}

func SendMoney(c *gin.Context) {
	var wal_id int
	wal_id_str := c.Param("walletId")
	wal_id, err := strconv.Atoi(wal_id_str)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	func(id1 int, id2 int) {
		// send money
	}(wal_id, req.To)

	c.Status(200)

	return
}

func GetWalletHistory(c *gin.Context) {
	wal_id := c.Param("walletId")
	c.String(http.StatusOK, "History of %s", wal_id)
}

func GetWalletState(c *gin.Context) {
	wal_id := c.Param("walletId")
	c.String(http.StatusOK, "Balance of %s is %d", wal_id, 5)
}
