package service

import (
	"time"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/messages"
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
	err = ws.Repo.TransferBalance(senderId, receiverId, amount)
	return
}

func (ws *WalletService) GetHistory(walId string) (history messages.History, err error) {
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
