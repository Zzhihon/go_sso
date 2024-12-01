package domain

import (
	"github.com/Zzhihon/sso/errs"
	"github.com/Zzhihon/sso/utils"
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
func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.token.SignedString([]byte(utils.SECRET))
	if err != nil {
		log.Println("Failed to sign access token " + err.Error())
		return "", errs.NewUnexpectedError("Failed to sign access token " + err.Error())
	}
	return signedString, nil
}

func (t AuthToken) newRefreshToken() (string, *errs.AppError) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshTokenClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedString, err := token.SignedString([]byte(utils.SECRET))
	if err != nil {
		log.Println("Failed to sign refresh token" + err.Error())
		return "", errs.NewUnexpectedError("Failed to sign refresh token " + err.Error())
	}
	return signedString, nil
}

func NewAccessTokenFromRefreshToken(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.SECRET), nil
	})
	if err != nil {
		log.Println("Failed to parse refresh token" + err.Error())
		return "", errs.NewUnexpectedError("Failed to parse refresh token " + err.Error())
	}
	r := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := r.AccessTokenClaims()
	authToken := NewAuthToken(accessTokenClaims)

	return authToken.NewAccessToken()
}
