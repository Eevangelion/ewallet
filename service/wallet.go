package service

import (
	"github.com/Eevangelion/ewallet/contracts"
)

type IWalletService interface {
	Create(balance float32) contracts.WalletResponse
	BalanceTransfer(sender_id int, receiver_id int, amount float32) error
	GetWalletHistory(wal_id int) error
	GetWalletState(wal_id int) error
}

type WalletService struct {
}

func (ws *WalletService) Create(balance float32) contracts.WalletResponse {
	return contracts.WalletResponse{
		Id:      1,
		Balance: 5.55,
	}
}

func (ws *WalletService) BalanceTransfer(sender_id int, receiver_id int, amount float32) (err error) {
	return err
}

func (ws *WalletService) GetWalletHistory(wal_id int) (err error) {
	return err
}

func (ws *WalletService) GetWalletState(wal_id int) (err error) {
	return err
}
