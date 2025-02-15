package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
)

var (
	ErrInvalidChars  = errors.New("contains invalid characters")
	ErrNeedMoreChars = errors.New("must be more than 8 characters")
	ErrNoUser            = errors.New("user not found")
	ErrWrongPass = errors.New("wrong email or password")
)

func WriteErrorResponse(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(entities.ErrorResponse{
		Errors: err.Error(),
	})
}
