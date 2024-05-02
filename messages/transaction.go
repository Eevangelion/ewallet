package messages

import "github.com/Eevangelion/ewallet/contracts"

type History struct {
	TransactionList []*contracts.Transaction
}
