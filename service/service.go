package service

import (
	"github.com/Eevangelion/ewallet/db/repositories"
)

type WalletService struct {
	Repo repositories.IWalletRepository
}

var walletService *WalletService

func GetWalletService(repo repositories.IWalletRepository) *WalletService {
	if walletService == nil {
		walletService = &WalletService{Repo: repo}
	}
	return walletService
}
