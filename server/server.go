package server

import "github.com/Eevangelion/ewallet/service"

type WalletServer struct {
	service *service.WalletService
}

func NewWalletServer(service *service.WalletService) *WalletServer {
	return &WalletServer{service: service}
}
