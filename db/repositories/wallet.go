package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/Eevangelion/ewallet/db"
	"github.com/Eevangelion/ewallet/logger"
	"github.com/Eevangelion/ewallet/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type WalletRepository struct {
}

const createWallet = `
INSERT INTO public."wallet" 
VALUES ($1, $2) 
RETURNING id`

func (wr *WalletRepository) Create(balance float32) (id string, err error) {
	logger := logger.GetLogger()
	pool, err := db.GetPool()
	if err != nil {
		logger.Error(
			"Error while connecting to DB:",
			zap.String("event", "connect_database"),
			zap.String("error", err.Error()),
		)
		return
	}
	err = pool.QueryRow(context.TODO(), createWallet, uuid.New(), balance).Scan(&id)
	if err != nil {
		logger.Error(
			"Error while creating wallet:",
			zap.String("event", "create_wallet"),
			zap.String("error", err.Error()),
		)
		return
	}
	return
}

const getBalance = `
SELECT balance
FROM public."wallet"
WHERE id = ($1)`

func (wr *WalletRepository) GetBalance(id string) (balance float32, err error) {
	logger := logger.GetLogger()
	pool, err := db.GetPool()
	if err != nil {
		logger.Error(
			"Error while connecting to DB:",
			zap.String("event", "connect_database"),
			zap.String("error", err.Error()),
		)
		return
	}
	err = pool.QueryRow(context.TODO(), getBalance, id).Scan(&balance)
	if err != nil {
		logger.Error(
			"Error while getting wallet balance:",
			zap.String("event", "get_balance"),
			zap.String("error", err.Error()),
		)
		return
	}
	return
}

const changeBalance = `
UPDATE public."wallet"
SET balance = ($1)
WHERE id = ($2)`

const createTransaction = `
INSERT INTO public."transaction"(sender_id, receiver_id, amount, time_stamp)
VALUES ($1, $2, $3, $4)`

func (wr *WalletRepository) TransferBalance(senderId string, receiverId string, amount float32) (err error) {
	logger := logger.GetLogger()
	pool, err := db.GetPool()
	if err != nil {
		logger.Error(
			"Error while connecting to DB:",
			zap.String("event", "connect_database"),
			zap.String("error", err.Error()),
		)
		return
	}
	tx, err := pool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		logger.Error(
			"Error while starting transaction:",
			zap.String("event", "transaction_start"),
			zap.String("error", err.Error()),
		)
		return
	}
	defer tx.Rollback(context.TODO())
	var senderBalance float32
	err = pool.QueryRow(context.TODO(), getBalance, senderId).Scan(&senderBalance)

	if err != nil {
		logger.Error(
			"Error while getting wallet balance:",
			zap.String("event", "get_balance"),
			zap.String("error", err.Error()),
		)
		return
	}

	if senderBalance < amount {
		err = errors.New("sender balance is less than transfer amount")
		logger.Error(
			"Error validating transfer amount:",
			zap.String("event", "validate_amount"),
			zap.String("error", err.Error()),
		)
		return
	}

	var receiverBalance float32
	err = pool.QueryRow(context.TODO(), getBalance, receiverId).Scan(&receiverBalance)
	if err != nil {
		logger.Error(
			"Error while getting wallet balance:",
			zap.String("event", "get_balance"),
			zap.String("error", err.Error()),
		)
		return
	}

	_, err = pool.Exec(context.TODO(), changeBalance, senderBalance-amount, senderId)

	if err != nil {
		logger.Error(
			"Error while changing wallet balance balance:",
			zap.String("event", "change_balance"),
			zap.String("error", err.Error()),
		)
		return
	}

	_, err = pool.Exec(context.TODO(), changeBalance, receiverBalance+amount, receiverId)

	if err != nil {
		logger.Error(
			"Error while changing wallet balance:",
			zap.String("event", "change_balance"),
			zap.String("error", err.Error()),
		)
		return
	}

	_, err = pool.Query(context.TODO(), createTransaction, senderId, receiverId, amount, time.Now())

	if err != nil {
		logger.Error(
			"Error while creating transfer transaction:",
			zap.String("event", "create_transfer_transaction"),
			zap.String("error", err.Error()),
		)
		return
	}

	err = tx.Commit(context.TODO())
	return
}

const getHistory = `
SELECT sender_id, receiver_id, amount, time_stamp
FROM public."transaction"
WHERE sender_id = ($1) OR receiver_id = ($1)`

func (wr *WalletRepository) GetHistory(id string) (txns []*models.Transcation, err error) {
	logger := logger.GetLogger()
	pool, err := db.GetPool()
	if err != nil {
		logger.Error(
			"Error while connecting to DB:",
			zap.String("event", "connect_database"),
			zap.String("error", err.Error()),
		)
		return
	}
	rows, err := pool.Query(context.TODO(), getHistory, id)

	if err != nil {
		logger.Error(
			"Error while getting history:",
			zap.String("event", "get_transfer_history"),
			zap.String("error", err.Error()),
		)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var txn models.Transcation
		err = rows.Scan(&txn.SenderId, &txn.ReceiverId, &txn.Amount, &txn.Timestamp)
		if err != nil {
			logger.Error(
				"Error while scanning rows:",
				zap.String("event", "scan_rows"),
				zap.String("error", err.Error()),
			)
			return
		}
		txns = append(txns, &txn)
	}

	if err = rows.Err(); err != nil {
		logger.Error(
			"Error while reading rows:",
			zap.String("event", "read_rows"),
			zap.String("error", err.Error()),
		)
		return
	}

	return
}
