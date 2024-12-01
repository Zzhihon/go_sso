package domain

import (
	"github.com/dgrijalva/jwt-go"
	"log"
)

type AuthToken struct {
	token *jwt.Token
}

// 获得meta token
func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}

// 生成access token
func (t AuthToken) NewAccessToken() (string, error) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed to sign access token" + err.Error())
		return "", err
	}
	return signedString, nil
}

// 生成refresh token
func (t AuthToken) VerifyAccessToken() (string, error) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshTokenClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed to sign refresh token" + err.Error())
		return "", err
	}
	return signedString, nil
}

func (t AuthToken) newRefreshToken() (string, error) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshTokenClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed to sign refresh token" + err.Error())
		return "", err
	}
	return signedString, nil
}
