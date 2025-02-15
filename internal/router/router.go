package router

import (
	"net/http"

	"github.com/KonstantinGalanin/itemStore/internal/handlers"
	"github.com/KonstantinGalanin/itemStore/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(userHandler *handlers.UserHandler) http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth", userHandler.Auth).Methods(http.MethodPost)

	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/info", userHandler.GetInfo).Methods(http.MethodGet)
	protected.HandleFunc("/sendCoin", userHandler.SendCoin).Methods(http.MethodPost)
	protected.HandleFunc("/buy/{item}", userHandler.BuyItem).Methods(http.MethodPost)
	
	return r
}
