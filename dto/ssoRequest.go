package dto

import (
	"errors"
	"github.com/Zzhihon/sso/utils"
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
	AccessToken  string `json:"access_token"`
	RefrestToken string `json:"refresh_token"`
}

func (r RefreshRequest) IsAccessTokenValid() *jwt.ValidationError {

	_, err := jwt.Parse(r.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.SECRET), nil
	})
	if err != nil {
		var vErr *jwt.ValidationError
		if errors.As(err, &vErr) {
			return vErr
		}
	}
	return nil
}
