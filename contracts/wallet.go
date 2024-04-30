package contracts

type Wallet struct {
	Balance float32 `json:"balance"`
}

type RequestSendMoney struct {
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

type WalletResponse struct {
	Id      string  `json:"id"`
	Balance float32 `json:"balance"`
}
