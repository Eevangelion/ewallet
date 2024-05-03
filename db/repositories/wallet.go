package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/Eevangelion/ewallet/db"
	"github.com/Eevangelion/ewallet/errs"
	"github.com/Eevangelion/ewallet/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletRepository struct {
}

const (
	SQLCreateWallet = `
	INSERT INTO public."wallet" 
	VALUES ($1, $2) 
	RETURNING id`

	SQLGetBalance = `
	SELECT balance
	FROM public."wallet"
	WHERE id = ($1)`

	SQLChangeBalance = `
	UPDATE public."wallet"
	SET balance = ($1)
	WHERE id = ($2)`

	SQLCreateTransaction = `
	INSERT INTO public."transaction"(sender_id, receiver_id, amount, time_stamp)
	VALUES ($1, $2, $3, $4)`

	SQLGetHistory = `
	SELECT sender_id, receiver_id, amount, time_stamp
	FROM public."transaction"
	WHERE sender_id = ($1) OR receiver_id = ($1)`

	SQLCheckWalletExists = `
	SELECT EXISTS(SELECT 1 FROM public."wallet" WHERE id=($1))`
)

func (wr *WalletRepository) Create(balance float32) (id string, e *errs.Err) {
	pool, err := db.GetPool()
	if err != nil {
		e = errs.NewErr(err, "connect_database", "")
		e = errs.WrapErr(e, "Create:")
		return
	}
	err = pool.QueryRow(context.TODO(), SQLCreateWallet, uuid.New(), balance).Scan(&id)
	if err != nil {
		e = errs.NewErr(err, "create_wallet", "")
		e = errs.WrapErr(e, "Create:")
		return
	}
	return
}

func (wr *WalletRepository) GetBalance(id string) (balance float32, err *errs.Err) {
	pool, e := db.GetPool()
	if e != nil {
		err = errs.NewErr(e, "connect_database", "")
		err = errs.WrapErr(err, "GetBalance:")
		return
	}
	e = pool.QueryRow(context.TODO(), SQLGetBalance, id).Scan(&balance)
	if e != nil {
		err = errs.NewErr(e, "get_balance", "not found")
		err = errs.WrapErr(err, "GetBalance:")
		return
	}
	return
}

func (wr *WalletRepository) TransferBalance(senderId string, receiverId string, amount float32) (err *errs.Err) {
	pool, e := db.GetPool()
	if err != nil {
		err = errs.NewErr(e, "connect_database", "")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}
	tx, e := pool.BeginTx(context.TODO(), pgx.TxOptions{})
	if e != nil {
		err = errs.NewErr(e, "transaction_start", "")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}
	defer tx.Rollback(context.TODO())
	var senderBalance float32
	e = pool.QueryRow(context.TODO(), SQLGetBalance, senderId).Scan(&senderBalance)

	if e != nil {
		err = errs.NewErr(e, "get_balance", "not found")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	if senderBalance < amount {
		e = errors.New("sender balance is less than transfer amount")
		err = errs.NewErr(e, "validate_amount", "bad data")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	var receiverBalance float32
	e = pool.QueryRow(context.TODO(), SQLGetBalance, receiverId).Scan(&receiverBalance)
	if e != nil {
		err = errs.NewErr(e, "get_balance", "not found")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	res, e := pool.Exec(context.TODO(), SQLChangeBalance, senderBalance-amount, senderId)

	if e != nil {
		err = errs.NewErr(e, "change_balance", "")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	if res.RowsAffected() == 0 {
		err = errs.NewErr(e, "change_balance", "not found")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	res, e = pool.Exec(context.TODO(), SQLChangeBalance, receiverBalance+amount, receiverId)

	if e != nil {
		err = errs.NewErr(e, "change_balance", "")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	if res.RowsAffected() == 0 {
		err = errs.NewErr(e, "change_balance", "not found")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	_, e = pool.Exec(context.TODO(), SQLCreateTransaction, senderId, receiverId, amount, time.Now())

	if e != nil {
		err = errs.NewErr(e, "create_transfer_transaction", "")
		err = errs.WrapErr(err, "TransferBalance:")
		return
	}

	e = tx.Commit(context.TODO())

	if e != nil {
		err = errs.NewErr(e, "commit_transaction", "")
		err = errs.WrapErr(err, "TransferBalance:")
	}
	return
}

func (wr *WalletRepository) GetHistory(id string) (txns []*models.Transcation, err *errs.Err) {
	pool, e := db.GetPool()
	if e != nil {
		err = errs.NewErr(e, "connect_database", "")
		err = errs.WrapErr(err, "GetHistory:")
		return
	}

	var exists bool
	e = pool.QueryRow(context.TODO(), SQLCheckWalletExists, id).Scan(&exists)

	if e != nil {
		err = errs.NewErr(e, "check_if_wallet_exists", "")
		err = errs.WrapErr(err, "GetHistory:")
		return
	}
	if !exists {
		e = errors.New("wallet not exists")
		err = errs.NewErr(e, "check_if_wallet_exists", "not found")
		err = errs.WrapErr(err, "GetHistory:")
		return
	}

	rows, e := pool.Query(context.TODO(), SQLGetHistory, id)

	if e != nil {
		err = errs.NewErr(e, "get_transfer_history", "")
		err = errs.WrapErr(err, "GetHistory:")
		return
	}

	defer rows.Close()

	for rows.Next() {
		var txn models.Transcation
		e = rows.Scan(&txn.SenderId, &txn.ReceiverId, &txn.Amount, &txn.Timestamp)
		if e != nil {
			err = errs.NewErr(e, "scan_rows", "")
			err = errs.WrapErr(err, "GetHistory:")
			return
		}
		txns = append(txns, &txn)
	}

	if e = rows.Err(); err != nil {
		err = errs.NewErr(e, "read_rows", "")
		err = errs.WrapErr(err, "GetHistory:")
		return
	}

	return
}
