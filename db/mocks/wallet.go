package mocks

import (
	"errors"
	"time"

	"github.com/Eevangelion/ewallet/errs"
	"github.com/Eevangelion/ewallet/models"
	"github.com/Eevangelion/ewallet/server"
	"github.com/google/uuid"
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

func (ms *WalletRepository) Create(balance float32) (id string, err *errs.Err) {
	uid := uuid.New()
	id = uid.String()
	ms.Wallets[id] = &models.Wallet{Balance: server.DefaultBalance}
	return
}

func (ms *WalletRepository) TransferBalance(senderId string, receiverId string, amount float32) (err *errs.Err) {
	if _, ok := ms.Wallets[senderId]; !ok {
		e := errors.New("wallet not exists")
		err = errs.NewErr(e, "check_if_wallet_exists", "not found")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}
	if _, ok := ms.Wallets[receiverId]; !ok {
		e := errors.New("wallet not exists")
		err = errs.NewErr(e, "check_if_wallet_exists", "not found")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	if ms.Wallets[senderId].Balance < amount {
		e := errors.New("sender balance is less than transfer amount")
		err = errs.NewErr(e, "validate_amount", "bad data")
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

func (ms *WalletRepository) GetHistory(id string) (txns []*models.Transcation, err *errs.Err) {
	if _, ok := ms.Wallets[id]; !ok {
		e := errors.New("wallet not exists")
		err = errs.NewErr(e, "check_if_wallet_exists", "not found")
		err = errs.WrapErr(err, "GetHistory:")
		return
	}
	for _, txn := range ms.Transactions {
		if txn.ReceiverId == id || txn.SenderId == id {
			txns = append(txns, txn)
		}
	}
	return
}

func (ms *WalletRepository) GetBalance(id string) (balance float32, err *errs.Err) {
	if _, ok := ms.Wallets[id]; !ok {
		e := errors.New("wallet not exists")
		err = errs.NewErr(e, "check_if_wallet_exists", "not found")
		err = errs.WrapErr(err, "GetBalance:")
		return
	}
	balance = ms.Wallets[id].Balance
	return
}
