package service

import (
	"errors"
	"time"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/logger"
	"github.com/Eevangelion/ewallet/messages"
	"github.com/Eevangelion/ewallet/utility"
	"go.uber.org/zap"
)

func (ws *WalletService) Create(balance float32) (wallet *messages.CreateWallet, err error) {
	id, err := ws.Repo.Create(balance)

	if err != nil {
		return
	}

	wallet = &messages.CreateWallet{
		Id:      id,
		Balance: balance,
	}

	return
}

func (ws *WalletService) BalanceTransfer(senderId string, receiverId string, amount float32) (err error) {
	logger := logger.GetLogger()
	if amount < 0 {
		err = errors.New("negative amount")
		logger.Error(
			"Error while getting transfering balance:",
			zap.String("event", "validate_amount"),
			zap.String("error", err.Error()),
		)
		return
	}
	if !utility.IsValidUUID(senderId) {
		err = errors.New("not valid wallet id")
		logger.Error(
			"Error while getting transfering balance:",
			zap.String("event", "validate_sender_id"),
			zap.String("error", err.Error()),
		)
		return
	}
	if !utility.IsValidUUID(receiverId) {
		err = errors.New("not valid wallet id")
		logger.Error(
			"Error while getting transfering balance:",
			zap.String("event", "validate_receiver_id"),
			zap.String("error", err.Error()),
		)
		return
	}
	err = ws.Repo.TransferBalance(senderId, receiverId, amount)
	return
}

func (ws *WalletService) GetHistory(walId string) (history messages.History, err error) {
	logger := logger.GetLogger()
	if !utility.IsValidUUID(walId) {
		err = errors.New("not valid wallet id")
		logger.Error(
			"Error while getting wallet history:",
			zap.String("event", "validate_wallet_id"),
			zap.String("error", err.Error()),
		)
		return
	}

	txns, err := ws.Repo.GetHistory(walId)

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

func (ws *WalletService) GetState(walId string) (wallet *messages.GetWalletState, err error) {
	logger := logger.GetLogger()
	if !utility.IsValidUUID(walId) {
		err = errors.New("not valid wallet id")
		logger.Error(
			"Error while getting wallet history:",
			zap.String("event", "validate_wallet_id"),
			zap.String("error", err.Error()),
		)
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
