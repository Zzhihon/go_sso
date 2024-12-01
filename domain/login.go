package domain

import (
	"database/sql"
	"github.com/Zzhihon/sso/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Login struct {
	UserID   string         `json:"userID" db:"username"`
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
			ExpiresAt: time.Now().Add(utils.ACCESS_TOKEN_DURATION).Unix()}}
}
