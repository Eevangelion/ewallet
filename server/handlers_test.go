package server_test

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/models"
	"github.com/Eevangelion/ewallet/server"
	"github.com/gin-gonic/gin"
)

type MockWalletService struct {
	Wallets map[string]models.Wallet
}

func (ms *MockWalletService) Create(balance float32) (*contracts.WalletResponse, error) {
	return &contracts.WalletResponse{
		Id:      "lol",
		Balance: balance,
	}, nil
}

func (ms *MockWalletService) BalanceTransfer(sender_id string, receiver_id string, amount float32) (err error) {
	return err
}

func (ms *MockWalletService) GetWalletHistory(wal_id string) (err error) {
	if _, ok := ms.Wallets[wal_id]; ok {
		log.Print("Return wallet history")
		// return wallet history
	} else {
		err = errors.New("Can't find wallet")
	}
	return err
}

func (ms *MockWalletService) GetWalletState(wal_id string) (err error) {
	if _, ok := ms.Wallets[wal_id]; ok {
		log.Print("Return wallet state")
		// return wallet state
	} else {
		err = errors.New("Can't find wallet")
	}
	return err
}

func TestCreateWallet(t *testing.T) {
	t.Run("can create valid wallet", func(t *testing.T) {
		expectedWallet := contracts.WalletResponse{Id: "lol", Balance: server.DefaultBalance}

		service := &MockWalletService{}
		server := server.NewWalletServer(service)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = req

		server.CreateWallet(c)

		assertStatus(t, w.Code, http.StatusOK)

		var wallet contracts.WalletResponse

		err := json.NewDecoder(w.Body).Decode(&wallet)
		if err != nil {
			t.Errorf("parse body wallet error")
		}
		req.Body.Close()

		if wallet.Id != expectedWallet.Id {
			t.Errorf("expected id %s but got %s", expectedWallet.Id, wallet.Id)
		}

		if wallet.Balance != expectedWallet.Balance {
			t.Errorf("expected balance %f but got %f", expectedWallet.Balance, wallet.Balance)
		}

	})
}

func assertStatus(t *testing.T, got int, want int) {
	if got != want {
		t.Errorf("wanted http status %d but got %d", got, want)
	}
}

func IsWalletValid(t *testing.T, wallet models.Wallet) {
	if wallet.Balance < 0 {
		t.Errorf("negative balance")
	}
}
