package jwt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
	jwtToken "github.com/golang-jwt/jwt/v5"
)

const (
	ExpTime = 7 * 24 * time.Hour
)

var (
	TokenSecret = []byte("secret")
)

type JWTInfo struct {
	Username string `json:"username"`
	jwtToken.RegisteredClaims
}

type JwtService struct{}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (j *JwtService) CreateToken(userItem *entities.User) ([]byte, error) {
	claims := JWTInfo{
		Username: userItem.Username,
		RegisteredClaims: jwtToken.RegisteredClaims{
			IssuedAt:  jwtToken.NewNumericDate(time.Now()),
			ExpiresAt: jwtToken.NewNumericDate(time.Now().Add(ExpTime)),
		},
	}

	token := jwtToken.NewWithClaims(jwtToken.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(TokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token response: %v", err)
	}

	return resp, nil
}

func GetToken(tokenString string) (string, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &JWTInfo{}
	token, err := jwtToken.ParseWithClaims(tokenString, claims, func(t *jwtToken.Token) (interface{}, error) {
		return TokenSecret, nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.Username, nil
}