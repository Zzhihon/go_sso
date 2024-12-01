package dto

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type LoginRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type RefreshRequest struct {
	AccessToken  string `json:"accessToken"`
	RefrestToken string `json:"refreshToken"`
}

func (r RefreshRequest) IsAccessTokenValid() *jwt.ValidationError {
	_, err := jwt.Parse(r.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		var vErr *jwt.ValidationError
		if errors.As(err, &vErr) {
			return vErr
		}
	}
	return nil
}
