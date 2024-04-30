package service

import (
	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/db"
)

type IWalletService interface {
	Create(balance float32) (*contracts.WalletResponse, error)
	BalanceTransfer(sender_id string, receiver_id string, amount float32) error
	GetWalletHistory(wal_id string) error
	GetWalletState(wal_id string) error
}

type WalletService struct {
}

func (ws *WalletService) Create(balance float32) (wallet *contracts.WalletResponse, err error) {
	id, err := db.Create(balance)

	if err != nil {
		return
	}

	wallet = &contracts.WalletResponse{
		Id:      id,
		Balance: balance,
	}

	return
}

func (ws *WalletService) BalanceTransfer(sender_id string, receiver_id string, amount float32) (err error) {
	return err
}

func (ws *WalletService) GetWalletHistory(wal_id string) (err error) {
	return err
}

func (ws *WalletService) GetWalletState(wal_id string) (err error) {
	return err
}
