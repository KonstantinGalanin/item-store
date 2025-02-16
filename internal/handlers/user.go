package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
	"github.com/KonstantinGalanin/itemStore/internal/utils"
	"github.com/gorilla/mux"
)

var (
	usernameValid = regexp.MustCompile(`[a-zA-Z0-9]+`)
	passwordValid = regexp.MustCompile(`.{8,}`)
)

func Validate(username, password string) error {
	if !usernameValid.MatchString(username) {
		return utils.ErrInvalidChars
	}

	if !passwordValid.MatchString(password) {
		return utils.ErrNeedMoreChars
	}

	return nil
}


type JwtService interface {
	CreateToken(userItem *entities.User) ([]byte, error)
}

//go:generate mockgen -source=user.go -destination=../service/user_service_mock.go -package=service
type UserService interface {
	BuyItem(userName string, itemName string) error
	SendCoin(fromUser, toUser string, amount int) error
	GetInfo(userName string) (*entities.InfoResponse, error)
	Auth(userName, password string) (*entities.User, error)
}

type UserHandler struct {
	UserService UserService
	JwtService  JwtService
}

func NewUserHandler(userService UserService, jwtService  JwtService) *UserHandler {
	return &UserHandler{
		UserService: userService,
		JwtService: jwtService,
	}
}

func (u *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := Validate(data.Username, data.Password); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := u.UserService.Auth(data.Username, data.Password)
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusUnauthorized)
		return 
	}

	resp, err := u.JwtService.CreateToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *UserHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	userName, ok := r.Context().Value("user").(string)
	if !ok {
		utils.WriteErrorResponse(w, fmt.Errorf("User not found"), http.StatusUnauthorized)
		return
	}

	var data struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := u.UserService.SendCoin(userName, data.ToUser, data.Amount); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	userName, ok := r.Context().Value("user").(string)
	if !ok {
		utils.WriteErrorResponse(w, fmt.Errorf("User not found"), http.StatusUnauthorized)
		return
	}

	itemName, exists := mux.Vars(r)["item"]
	if !exists {
		utils.WriteErrorResponse(w, fmt.Errorf("Item not exists"), http.StatusBadRequest)
		return
	}

	if err := u.UserService.BuyItem(userName, itemName); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	userName, ok := r.Context().Value("user").(string)
	if !ok {
		utils.WriteErrorResponse(w, fmt.Errorf("User not found"), http.StatusUnauthorized)
		return
	}

	info, err := u.UserService.GetInfo(userName)
	if err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		utils.WriteErrorResponse(w, err, http.StatusInternalServerError)
	}
}
