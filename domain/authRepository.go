package domain

import "github.com/Zzhihon/sso/errs"

type AuthRepository interface {
	FindBy(userID string) (*Login, *errs.AppError)
	GenerateRefreshToken(token AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
}
