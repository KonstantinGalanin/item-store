package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KonstantinGalanin/itemStore/internal/handlers"
	"github.com/KonstantinGalanin/itemStore/internal/middleware"
	repository "github.com/KonstantinGalanin/itemStore/internal/repository/user"
	"github.com/KonstantinGalanin/itemStore/internal/service"
	"github.com/KonstantinGalanin/itemStore/pkg/jwt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestSendCoin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserPostgresRepo(db)
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService, jwt.NewJwtService())

	router := mux.NewRouter()
	router.HandleFunc("/api/auth", userHandler.Auth).Methods("POST")
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/api/sendCoin", userHandler.SendCoin).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	var data struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	data.ToUser = Receiver
	data.Amount = 100

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal auth data: %v", err)
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/sendCoin", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	token := authenticateAndGetToken(t, ts, Sender, "pass1234")
	req.Header.Set("Authorization", "Bearer "+token) //

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var senderBalance, receiverBalance int
	err = db.QueryRow("SELECT balance FROM users WHERE username = $1", Sender).Scan(&senderBalance)
	if err != nil {
		t.Fatalf("Failed to get sender balance: %v", err)
	}
	assert.Equal(t, 900, senderBalance)

	err = db.QueryRow("SELECT balance FROM users WHERE username = $1", Receiver).Scan(&receiverBalance)
	if err != nil {
		t.Fatalf("Failed to get receiver balance: %v", err)
	}
	assert.Equal(t, 1100, receiverBalance)

	var exchangeAmount int
	err = db.QueryRow("SELECT amount FROM exchanges WHERE from_id = (SELECT id FROM users WHERE username = $1) AND to_id = (SELECT id FROM users WHERE username = $2)", Sender, Receiver).Scan(&exchangeAmount)
	if err != nil {
		t.Fatalf("Failed to get exchange record: %v", err)
	}
	assert.Equal(t, 100, exchangeAmount)
}
