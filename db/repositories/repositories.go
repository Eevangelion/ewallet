package repositories

import (
	"github.com/Eevangelion/ewallet/models"
)

type IWalletRepository interface {
	Create(float32) (string, error)
	GetBalance(string) (float32, error)
	GetHistory(string) ([]*models.Transcation, error)
	TransferBalance(string, string, float32) error
}

var walletRepo *WalletRepository

func GetWalletRepository() *WalletRepository {
	if walletRepo == nil {
		walletRepo = &WalletRepository{}
	}
	return walletRepo
}
