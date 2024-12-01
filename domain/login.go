package domain

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const TOKEN_DURATION = time.Hour * 24

type Login struct {
	UserID   string         `json:"userID" db:"userID"`
	Password string         `json:"password" db:"password"`
	Name     string         `json:"name" db:"name"`
	Role     sql.NullString `json:"role" db:"role"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	return AccessTokenClaims{
		UserID: l.UserID,
		Name:   l.Name,
		Role:   l.Role.String,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix()}}
}
