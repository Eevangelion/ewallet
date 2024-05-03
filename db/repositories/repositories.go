package repositories

import (
	"github.com/Eevangelion/ewallet/errs"
	"github.com/Eevangelion/ewallet/models"
)

type IWalletRepository interface {
	Create(float32) (string, *errs.Err)
	GetBalance(string) (float32, *errs.Err)
	GetHistory(string) ([]*models.Transcation, *errs.Err)
	TransferBalance(string, string, float32) *errs.Err
}

var walletRepo *WalletRepository

func GetWalletRepository() *WalletRepository {
	if walletRepo == nil {
		walletRepo = &WalletRepository{}
	}
	return walletRepo
}
