package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Eevangelion/ewallet/contracts"
	"github.com/Eevangelion/ewallet/db/mocks"
	"github.com/Eevangelion/ewallet/server"
	"github.com/Eevangelion/ewallet/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestCreateWallet(t *testing.T) {
	t.Run("can create valid wallet", func(t *testing.T) {

		repo := mocks.GetWalletRepository()
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = req

		serv.CreateWallet(c)

		assertStatus(t, w.Code, http.StatusOK)

		var wallet *contracts.WalletResponse

		err := json.NewDecoder(w.Body).Decode(&wallet)
		if err != nil {
			t.Errorf("parse body wallet error")
		}
		req.Body.Close()

		IsWalletValid(t, wallet)

		if wallet.Balance != server.DefaultBalance {
			t.Errorf("expected default balance %f but got %f", server.DefaultBalance, wallet.Balance)
		}

	})
}

func TestSendMoney(t *testing.T) {
	t.Run("can't transfer more than have", func(t *testing.T) {
		repo := mocks.GetWalletRepository()
		senderId, _ := repo.Create(server.DefaultBalance)
		receiverId, _ := repo.Create(server.DefaultBalance)
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)

		body := strings.NewReader(fmt.Sprintf(`{"to": %s, "amount": %d }`, receiverId, 110))

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/wallet/%s/send", senderId), body)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = req

		serv.SendMoney(c)

		assertStatus(t, w.Code, http.StatusBadRequest)
	})
}

func TestGetWalletHistory(t *testing.T) {
	t.Run("get transaction history", func(t *testing.T) {
		repo := mocks.GetWalletRepository()
		senderId, _ := repo.Create(server.DefaultBalance)
		receiverId, _ := repo.Create(server.DefaultBalance)
		repo.TransferBalance(senderId, receiverId, 50)
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)

		body := strings.NewReader(fmt.Sprintf(`{"to": %s, "amount": %d }`, receiverId, 110))

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/wallet/%s/history", senderId), body)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = req

		serv.GetWalletHistory(c)

		var txnHistory contracts.TransactionHistory

		err := json.NewDecoder(w.Body).Decode(&txnHistory.TransactionList)
		if err != nil {
			t.Errorf("parse body wallet history error")
		}

		if len(txnHistory.TransactionList) == 0 {
			t.Errorf("empty wallet history: expected len %d, got %d", 1, 0)
		}

		assertStatus(t, w.Code, http.StatusOK)
	})
}

func assertStatus(t *testing.T, got int, want int) {
	if got != want {
		t.Errorf("wanted http status %d but got %d", want, got)
	}
}

func IsWalletValid(t *testing.T, wallet *contracts.WalletResponse) {
	if wallet.Balance < 0 {
		t.Errorf("negative balance")
	}
	if err := uuid.Validate(wallet.Id); err != nil {
		t.Errorf("failed validation of wallet id, got %s", wallet.Id)
	}
}
