package integration

import (
	"bytes"
	"database/sql"
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

const (
	TestUser = "Konstantin"
	Receiver = "Masha"
	Sender   = "Dima"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "host=localhost port=5432 user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", TestUser, "pass1234")
	db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", Receiver, "pass1234")
	db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", Sender, "pass1234")

	return db
}

func authenticateAndGetToken(t *testing.T, ts *httptest.Server, username, password string) string {
	authData := map[string]string{
		"username": username,
		"password": password,
	}
	authBody, err := json.Marshal(authData)
	if err != nil {
		t.Fatalf("Failed to marshal auth data: %v", err)
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/auth", bytes.NewBuffer(authBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	var authResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("Failed to decode auth response: %v", err)
	}

	token, ok := authResponse["token"]
	if !ok {
		t.Fatalf("Token not found in auth response")
	}

	return token
}

func TestBuyItem(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserPostgresRepo(db)
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService, jwt.NewJwtService())

	router := mux.NewRouter()
	router.HandleFunc("/api/auth", userHandler.Auth).Methods("POST")
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/api/buy/{item}", userHandler.BuyItem).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	token := authenticateAndGetToken(t, ts, TestUser, "pass1234")

	req, err := http.NewRequest("POST", ts.URL+"/api/buy/t-shirt", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IktvbnN0YW50aW4iLCJleHAiOjE3NDAzMzQ3NDUsImlhdCI6MTczOTcyOTk0NX0.Ojz2HB4ffHtrVdBdpDRpHpL0qs3v7BEF5ztwHBIIGfY"
	req.Header.Set("Authorization", "Bearer "+token) //

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var balance int
	err = db.QueryRow("SELECT balance FROM users WHERE username = $1", TestUser).Scan(&balance)
	if err != nil {
		t.Fatalf("Failed to get user balance: %v", err)
	}
	assert.Equal(t, 920, balance)

	var itemID int
	err = db.QueryRow("SELECT item_id FROM purchases WHERE user_id = (SELECT id FROM users WHERE username = $1)", TestUser).Scan(&itemID)
	if err != nil {
		t.Fatalf("Failed to get item from purchases: %v", err)
	}
	assert.NotEqual(t, 0, itemID)
}
