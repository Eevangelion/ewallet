package service

import (
	"errors"
	"time"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/errs"
	"github.com/Eevangelion/ewallet/messages"
	"github.com/Eevangelion/ewallet/utility"
)

func (ws *WalletService) Create(balance float32) (wallet *messages.CreateWallet, err *errs.Err) {
	id, err := ws.Repo.Create(balance)

	if err != nil {
		err = errs.WrapErr(err, "Create")
		return
	}

	wallet = &messages.CreateWallet{
		Id:      id,
		Balance: balance,
	}

	return
}

func (ws *WalletService) BalanceTransfer(senderId string, receiverId string, amount float32) (err *errs.Err) {
	if amount <= 0 {
		e := errors.New("negative amount")
		err = errs.NewErr(e, "validate_amount", "bad data")
		err = errs.WrapErr(err, "BalanceTransfer:")
		return
	}
	if !utility.IsValidUUID(senderId) {
		e := errors.New("not valid wallet id")
		err = errs.NewErr(e, "validate_sender_id", "bad data")
		err = errs.WrapErr(err, "BalanceTransfer:")
		return
	}
	if !utility.IsValidUUID(receiverId) {
		e := errors.New("not valid wallet id")
		err = errs.NewErr(e, "validate_receiver_id", "bad data")
		err = errs.WrapErr(err, "BalanceTransfer:")
		return
	}
	err = ws.Repo.TransferBalance(senderId, receiverId, amount)
	if err != nil {
		err = errs.WrapErr(err, "BalanceTransfer:")
	}
	return
}

func (ws *WalletService) GetHistory(walId string) (history messages.History, err *errs.Err) {
	if !utility.IsValidUUID(walId) {
		e := errors.New("not valid wallet id")
		err = errs.NewErr(e, "validate_wallet_id", "bad data")
		err = errs.WrapErr(err, "GetHistory:")
		return
	}

	txns, err := ws.Repo.GetHistory(walId)

	if err != nil {
		err = errs.WrapErr(err, "GetHistory:")
		return
	}

	for _, txn := range txns {
		history.TransactionList = append(history.TransactionList, &contracts.Transaction{
			SenderId:   txn.SenderId,
			ReceiverId: txn.ReceiverId,
			Amount:     txn.Amount,
			Timestamp:  txn.Timestamp.Format(time.RFC3339),
		})
	}
	return
}

func (ws *WalletService) GetState(walId string) (wallet *messages.GetWalletState, err *errs.Err) {
	if !utility.IsValidUUID(walId) {
		e := errors.New("not valid wallet id")
		err = errs.NewErr(e, "validate_wallet_id", "bad data")
		err = errs.WrapErr(err, "GetState:")
		return
	}
	balance, err := ws.Repo.GetBalance(walId)

	if err != nil {
		return
	}

	wallet = &messages.GetWalletState{
		Id:      walId,
		Balance: balance,
	}
	return
}
