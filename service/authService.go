package service

import (
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/dto"
	"github.com/dgrijalva/jwt-go"
	"log"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, error)
	Verify(token string) (bool, error)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*string, error) {
	login, err := s.repo.FindBy(req.UserID, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s DefaultAuthService) Verify(token string) (bool, error) {
	jwtToken, err := jwtTokenFromString(token)
	if err != nil {
		return false, err
	} else {
		if jwtToken.Valid {
			return true, nil
		} else {
			return false, err
		}
	}
}

func jwtTokenFromString(tokenstring string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token" + err.Error())
		return nil, err
	}
	return token, nil
}

func NewAuthService(repo domain.AuthRepository) DefaultAuthService {
	return DefaultAuthService{repo: repo}
}
