package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
	"github.com/KonstantinGalanin/itemStore/internal/service"
	"github.com/KonstantinGalanin/itemStore/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	password = "password"
	id       = "1"
	username = "username"
)

var (
	FakeError = errors.New("fake writer error")
)

func TestValidate(t *testing.T) {
	type user struct {
		username string
		password string
		expected error
	}
	cases := []user{
		{
			username: "User",
			password: password,
			expected: nil,
		},
		{
			username: "User",
			password: "p",
			expected: utils.ErrNeedMoreChars,
		},
		{
			username: "<><{};';",
			password: password,
			expected: utils.ErrInvalidChars,
		},
	}

	for _, c := range cases {
		err := Validate(c.username, c.password)
		assert.Equal(t, c.expected, err)
	}
}

func TestNewUserPostgresRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)
	mockJwtService := service.NewMockJwtService(ctrl)
	handler := NewUserHandler(mockUserService, mockJwtService)

	assert.NotNil(t, handler)
}

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)
	mockJwtService := service.NewMockJwtService(ctrl)

	userHandler := UserHandler{
		UserService: mockUserService,
		JwtService:  mockJwtService,
	}

	t.Run("json parse error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer([]byte("{invalid json}")))
		w := httptest.NewRecorder()

		userHandler.Auth(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		data := map[string]string{"username": "", "password": ""}
		body, _ := json.Marshal(data)

		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		userHandler.Auth(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("auth error", func(t *testing.T) {
		data := map[string]string{"username": "testuser", "password": "wrongpass"}
		body, _ := json.Marshal(data)

		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockUserService.EXPECT().
			Auth("testuser", "wrongpass").
			Return(nil, errors.New("invalid credentials"))

		userHandler.Auth(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		data := map[string]string{"username": "testuser", "password": "correctpass"}
		body, _ := json.Marshal(data)

		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockUser := &entities.User{ID: 1, Username: "testuser"}
		mockToken := []byte(`{"token":"mocked_jwt"}`)

		mockUserService.EXPECT().
			Auth("testuser", "correctpass").
			Return(mockUser, nil)

		mockJwtService.EXPECT().
			CreateToken(mockUser).
			Return(mockToken, nil)

		userHandler.Auth(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody map[string]string
		json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.Equal(t, "mocked_jwt", responseBody["token"])
	})

	t.Run("jwt create token error", func(t *testing.T) {
		data := map[string]string{"username": "testuser", "password": "correctpass"}
		body, _ := json.Marshal(data)

		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockUser := &entities.User{ID: 1, Username: "testuser"}

		mockUserService.EXPECT().
			Auth("testuser", "correctpass").
			Return(mockUser, nil)

		mockJwtService.EXPECT().
			CreateToken(mockUser).
			Return(nil, errors.New("failed to marshal token response"))

		userHandler.Auth(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)

	userHandler := UserHandler{
		UserService: mockUserService,
	}

	t.Run("user context error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sendCoin", nil)
		w := httptest.NewRecorder()

		userHandler.SendCoin(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("json decode error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer([]byte("{invalid json}")))
		ctx := context.WithValue(req.Context(), "user", "alice")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		userHandler.SendCoin(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("send coin error", func(t *testing.T) {
		data := map[string]interface{}{
			"toUser": "bob",
			"amount": 100,
		}
		body, _ := json.Marshal(data)

		req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), "user", "alice")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserService.EXPECT().
			SendCoin("alice", "bob", 100).
			Return(errors.New("transaction failed"))

		userHandler.SendCoin(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		data := map[string]interface{}{
			"toUser": "bob",
			"amount": 50,
		}
		body, _ := json.Marshal(data)

		req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(body))
		ctx := context.WithValue(req.Context(), "user", "alice")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockUserService.EXPECT().
			SendCoin("alice", "bob", 50).
			Return(nil)

		userHandler.SendCoin(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestUserHandler_BuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)

	userHandler := UserHandler{
		UserService: mockUserService,
	}

	t.Run("error no user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/buy/cup", nil)
		w := httptest.NewRecorder()

		userHandler.BuyItem(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("wrong url", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/buy/", nil)
		ctx := context.WithValue(req.Context(), "user", "alice")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{})

		userHandler.BuyItem(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/buy/cup", nil)
		ctx := context.WithValue(req.Context(), "user", "alice")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"item": "cup"})

		mockUserService.EXPECT().
			BuyItem("alice", "cup").
			Return(errors.New("not enough coins"))

		userHandler.BuyItem(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/buy/cup", nil)
		ctx := context.WithValue(req.Context(), "user", "alice")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"item": "cup"})

		mockUserService.EXPECT().
			BuyItem("alice", "cup").
			Return(nil)

		userHandler.BuyItem(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestGetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)

	userHandler := UserHandler{
		UserService: mockUserService,
	}

	t.Run("success", func(t *testing.T) {
		userName := "testuser"
		userInfo := &entities.InfoResponse{
			Coins:     100,
			Inventory: make([]*entities.Item, 0),
			CoinHistory: entities.CoinHistory{
				Received: make([]*entities.ReceiveOperation, 0),
				Sent:     make([]*entities.SentOperation, 0),
			},
		}
		mockUserService.EXPECT().GetInfo(userName).Return(userInfo, nil)

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", userName))
		w := httptest.NewRecorder()

		userHandler.GetInfo(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}
