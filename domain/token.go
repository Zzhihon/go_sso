package domain

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
)

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Expiry   int64  `json:"exp"`
}

func BuildClaimsFromJwtMapClaims(mapClaims jwt.MapClaims) (*Claims, error) {
	bytes, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, err
	}
	var c Claims
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Claims) IsUserRole() bool {
	return c.Role == "user"
}

func (c Claims) IsAdminRole() bool {
	return c.Role == "admin"
}

func (c *Claims) IsValidName(name string) bool {
	return c.Username == name
}

func (c *Claims) IsValidID(id string) bool {
	return c.UserID == id
}

func (c Claims) IsRequestVerifyWithTokenClaims(urlParams map[string]string) bool {
	if !c.IsValidID(urlParams["userID"]) {
		return false
	}
	if !c.IsValidName(urlParams["username"]) {
		return false
	}

	return true
}
