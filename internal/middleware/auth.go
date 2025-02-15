package middleware

import (
	"context"
	"net/http"

	"github.com/KonstantinGalanin/itemStore/internal/utils"
	"github.com/KonstantinGalanin/itemStore/pkg/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := jwt.GetToken(r.Header.Get("Authorization"))
		if err != nil {
			utils.WriteErrorResponse(w, err, http.StatusUnauthorized) // or bad request
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
