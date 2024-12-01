package domain

import (
	"github.com/Zzhihon/sso/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type RefreshTokenClaims struct {
	TokenType string `json:"tokenType"`
	UserID    string `json:"userId"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

type AccessTokenClaims struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func (c AccessTokenClaims) RefreshTokenClaims() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType: "refreshToken",
		UserID:    c.UserID,
		Name:      c.Name,
		Role:      c.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(utils.REFRESH_TOKEN_DURATION).Unix(),
		},
	}
}

func (c RefreshTokenClaims) AccessTokenClaims() AccessTokenClaims {
	return AccessTokenClaims{
		UserID: c.UserID,
		Name:   c.Name,
		Role:   c.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(utils.ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}

func (c AccessTokenClaims) IsUserIDValid(id string) bool {
	return c.UserID == id
}

func (c AccessTokenClaims) IsUserRole() bool {
	return c.Role == "user"
}

func (c AccessTokenClaims) IsAdminRole() bool {
	return c.Role == "admin"
}

func (c AccessTokenClaims) IsNameValid(name string) bool {
	return c.Name != name
}
