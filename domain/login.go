package domain

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const TOKEN_DURATION = time.Hour * 24

type Login struct {
	UserID   string         `json:"userID" db:"userID"`
	Password string         `json:"password" db:"password"`
	Name     string         `json:"name" db:"name"`
	Role     sql.NullString `json:"role" db:"role"`
}

func (l Login) GenerateToken() (*string, error) {
	var claims jwt.MapClaims
	//这里看情况可以添加role.Valid的检验逻辑

	claims = l.MapClaimsForUSer()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed to sign token" + err.Error())
		return nil, err
	}
	return &signedTokenString, nil
}

func (l Login) MapClaimsForUSer() jwt.MapClaims {
	return jwt.MapClaims{
		"userID": l.UserID,
		"name":   l.Name,
		"role":   l.Role.String,
		"exp":    time.Now().Add(TOKEN_DURATION).Unix(),
	}
}
