package messages

type CreateWallet struct {
	Id      string
	Balance float32
}

type GetWalletState struct {
	Id      string
	Balance float32
}
