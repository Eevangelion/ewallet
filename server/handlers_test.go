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

		bodyJson := fmt.Sprintf(`{"to": "%s", "amount": %f }`, receiverId, server.DefaultBalance+1)

		body := strings.NewReader(bodyJson)

		url := fmt.Sprintf("/api/v1/wallet/%s/send", senderId)

		req := httptest.NewRequest(http.MethodPost, url, body)
		w := httptest.NewRecorder()

		c, e := gin.CreateTestContext(w)
		e.POST("/api/v1/wallet/:walletId/send", serv.SendMoney)
		c.Request = req
		e.ServeHTTP(w, req)

		assertStatus(t, w.Code, http.StatusBadRequest)
	})
	t.Run("can't transfer non positive amount", func(t *testing.T) {
		repo := mocks.GetWalletRepository()
		senderId, _ := repo.Create(server.DefaultBalance)
		receiverId, _ := repo.Create(server.DefaultBalance)
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)

		bodyJson := fmt.Sprintf(`{"to": "%s", "amount": %f }`, receiverId, -5.)

		body := strings.NewReader(bodyJson)

		url := fmt.Sprintf("/api/v1/wallet/%s/send", senderId)

		req := httptest.NewRequest(http.MethodPost, url, body)
		w := httptest.NewRecorder()

		c, e := gin.CreateTestContext(w)
		e.POST("/api/v1/wallet/:walletId/send", serv.SendMoney)
		c.Request = req
		e.ServeHTTP(w, req)

		assertStatus(t, w.Code, http.StatusBadRequest)
	})
}

func TestGetWalletHistory(t *testing.T) {
	t.Run("get transaction history", func(t *testing.T) {
		repo := mocks.GetWalletRepository()
		senderId, _ := repo.Create(server.DefaultBalance)
		receiverId, _ := repo.Create(server.DefaultBalance)
		repo.TransferBalance(senderId, receiverId, 1)
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallet/%s/history", senderId), nil)
		w := httptest.NewRecorder()

		c, e := gin.CreateTestContext(w)
		e.GET("/api/v1/wallet/:walletId/history", serv.GetWalletHistory)
		c.Request = req
		e.ServeHTTP(w, req)

		assertStatus(t, w.Code, http.StatusOK)

		var txnHistory contracts.TransactionHistory

		err := json.NewDecoder(w.Body).Decode(&txnHistory.TransactionList)
		if err != nil {
			t.Fatalf("parse body wallet history error")
		}

		if len(txnHistory.TransactionList) == 0 {
			t.Fatalf("empty wallet history: expected len %d, got %d", 1, 0)
		}
	})
	t.Run("don't create transaction after not valid transfer", func(t *testing.T) {
		repo := mocks.GetWalletRepository()
		senderId, _ := repo.Create(server.DefaultBalance)
		receiverId, _ := repo.Create(server.DefaultBalance)
		repo.TransferBalance(senderId, receiverId, server.DefaultBalance+1)
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallet/%s/history", senderId), nil)
		w := httptest.NewRecorder()

		c, e := gin.CreateTestContext(w)
		e.GET("/api/v1/wallet/:walletId/history", serv.GetWalletHistory)
		c.Request = req
		e.ServeHTTP(w, req)

		assertStatus(t, w.Code, http.StatusOK)

		var txnHistory contracts.TransactionHistory

		err := json.NewDecoder(w.Body).Decode(&txnHistory.TransactionList)
		if err != nil {
			t.Fatalf("parse body wallet history error")
		}

		if len(txnHistory.TransactionList) > 0 {
			t.Fatalf("not empty wallet history: expected len %d, got %d", 0, 1)
		}
	})

	t.Run("not found when wallet doesn't exist", func(t *testing.T) {
		repo := mocks.GetWalletRepository()
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)
		walId := uuid.New()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallet/%s/history", walId), nil)
		w := httptest.NewRecorder()

		c, e := gin.CreateTestContext(w)
		e.GET("/api/v1/wallet/:walletId/history", serv.GetWalletHistory)
		c.Request = req
		e.ServeHTTP(w, req)

		assertStatus(t, w.Code, http.StatusNotFound)
	})
}

func TestGetWalletState(t *testing.T) {
	t.Run("not found when wallet doesn't exist", func(t *testing.T) {
		repo := mocks.GetWalletRepository()
		svc := service.GetWalletService(repo)
		serv := server.NewWalletServer(svc)
		walId := uuid.New()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallet/%s", walId), nil)
		w := httptest.NewRecorder()

		c, e := gin.CreateTestContext(w)
		e.GET("/api/v1/wallet/:walletId", serv.GetWalletHistory)
		c.Request = req
		e.ServeHTTP(w, req)

		assertStatus(t, w.Code, http.StatusNotFound)
	})
}

func assertStatus(t *testing.T, got int, want int) {
	if got != want {
		t.Fatalf("wanted http status %d but got %d", want, got)
	}
}

func IsWalletValid(t *testing.T, wallet *contracts.WalletResponse) {
	if wallet.Balance < 0 {
		t.Fatalf("negative balance")
	}
	if err := uuid.Validate(wallet.Id); err != nil {
		t.Fatalf("failed validation of wallet id, got %s", wallet.Id)
	}
}
