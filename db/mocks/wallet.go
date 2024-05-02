package mocks

import (
	"errors"
	"time"

	"github.com/Eevangelion/ewallet/logger"
	"github.com/Eevangelion/ewallet/models"
	"github.com/Eevangelion/ewallet/server"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletRepository struct {
	Wallets      map[string]*models.Wallet
	Transactions []*models.Transcation
}

var walletRepo *WalletRepository

func GetWalletRepository() *WalletRepository {
	if walletRepo == nil {
		walletRepo = &WalletRepository{
			Wallets: make(map[string]*models.Wallet),
		}
	}
	walletRepo.Wallets = make(map[string]*models.Wallet)
	walletRepo.Transactions = nil
	return walletRepo
}

func (ms *WalletRepository) Create(balance float32) (id string, err error) {
	uid := uuid.New()
	id = uid.String()
	ms.Wallets[id] = &models.Wallet{Balance: server.DefaultBalance}
	return
}

func (ms *WalletRepository) TransferBalance(senderId string, receiverId string, amount float32) (err error) {
	logger := logger.GetLogger()
	if ms.Wallets[senderId].Balance < amount {
		err = errors.New("sender balance is less than transfer amount")
		logger.Error(
			"Error validating transfer amount:",
			zap.String("event", "validate_amount"),
			zap.String("error", err.Error()),
		)
		return
	}

	ms.Wallets[senderId].Balance -= amount
	ms.Wallets[receiverId].Balance += amount
	ms.Transactions = append(ms.Transactions, &models.Transcation{
		SenderId:   senderId,
		ReceiverId: receiverId,
		Amount:     amount,
		Timestamp:  time.Now(),
	})
	return
}

func (ms *WalletRepository) GetHistory(id string) (txns []*models.Transcation, err error) {
	txns = ms.Transactions
	return
}

func (ms *WalletRepository) GetBalance(id string) (balance float32, err error) {
	balance = ms.Wallets[id].Balance
	return
}
