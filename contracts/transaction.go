package contracts

type Transaction struct {
	Timestamp  string  `json:"time"`
	SenderId   string  `json:"from"`
	ReceiverId string  `json:"to"`
	Amount     float32 `json:"amount"`
}

type TransactionHistory struct {
	TransactionList []*Transaction
}
