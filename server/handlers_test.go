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
	Wallets []models.Wallet
}

func (ms *MockWalletService) Create(balance float32) contracts.WalletResponse {
	return contracts.WalletResponse{
		Id:      0,
		Balance: 100,
	}
}

func (ms *MockWalletService) BalanceTransfer(sender_id int, receiver_id int, amount float32) (err error) {
	return err
}

func (ms *MockWalletService) GetWalletHistory(wal_id int) (err error) {
	if wal_id >= len(ms.Wallets) {
		err = errors.New("Can't find wallet")
	} else {
		log.Print("Return wallet history")
		// return wallet history
	}
	return err
}

func (ms *MockWalletService) GetWalletState(wal_id int) (err error) {
	if wal_id >= len(ms.Wallets) {
		err = errors.New("Can't find wallet")
	} else {
		log.Print("Return wallet state")
		// return wallet history
	}
	return err
}

func TestCreateWallet(t *testing.T) {
	t.Run("can create valid wallet", func(t *testing.T) {
		expectedWallet := contracts.WalletResponse{Id: 0, Balance: server.DefaultBalance}

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
			t.Errorf("expected id %d but got %d", wallet.Id, expectedWallet.Id)
		}

		if wallet.Balance != expectedWallet.Balance {
			t.Errorf("expected balance %f but got %f", wallet.Balance, expectedWallet.Balance)
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
