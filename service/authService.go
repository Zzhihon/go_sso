package service

import (
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
	"github.com/Zzhihon/sso/utils"
	"github.com/dgrijalva/jwt-go"
	"log"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(token string) (bool, *errs.AppError)
	Refresh(request dto.RefreshRequest) (*dto.RefreshTokenResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo      domain.AuthRepository
	utilsRepo domain.UtilsRepository
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var login *domain.Login

	//检查用户名和密码是否正确
	_, pErr := s.utilsRepo.CheckPassword(req.UserID, req.Password)
	if pErr != nil {
		return nil, pErr
	}

	login, err := s.repo.FindBy(req.UserID)
	if err != nil {
		return nil, err
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)
	var accessToken, refreshToken string
	if accessToken, err = authToken.NewAccessToken(); err != nil {
		return nil, err
	}

	if refreshToken, err = s.repo.GenerateRefreshToken(authToken); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s DefaultAuthService) Verify(token string) (bool, *errs.AppError) {
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

func (s DefaultAuthService) Refresh(request dto.RefreshRequest) (*dto.RefreshTokenResponse, *errs.AppError) {
	if vErr := request.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {
			var appErr *errs.AppError
			if appErr = s.repo.RefreshTokenExists(request.RefrestToken); appErr != nil {
				return nil, appErr
			}

			var accessToken string
			var err *errs.AppError
			if accessToken, err = domain.NewAccessTokenFromRefreshToken(request.RefrestToken); err != nil {
				return nil, err
			}
			return &dto.RefreshTokenResponse{
				AccessToken: accessToken,
			}, nil
		}
		return nil, errs.NewUnAuthorizedError(vErr.Error())
	}
	return nil, errs.NewUnexpectedError("can not generate a new access token until the current one is expired")
}

func jwtTokenFromString(tokenstring string) (*jwt.Token, *errs.AppError) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token" + err.Error())
		return nil, errs.NewUnexpectedError(err.Error())
	}
	return token, nil
}

func NewAuthService(repo domain.AuthRepository, utilsRepo domain.UtilsRepository) DefaultAuthService {
	return DefaultAuthService{
		repo:      repo,
		utilsRepo: utilsRepo,
	}
}
