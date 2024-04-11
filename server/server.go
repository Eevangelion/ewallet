package server

import "github.com/Eevangelion/ewallet/service"

type WalletServer struct {
	service service.IWalletService
}

func NewWalletServer(service service.IWalletService) *WalletServer {
	return &WalletServer{service: service}
}
