package domain

import (
	"context"
	"github.com/Zzhihon/sso/errs"
	"github.com/go-redis/redis/v8"
	"time"
)

type AuthRepositoryRedis interface {
	GenerateRefreshToken(token AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
}

type authRepositoryRedisImpl struct {
	rdb *redis.Client
	ctx context.Context
}

func (a authRepositoryRedisImpl) GenerateRefreshToken(authtoken AuthToken) (string, *errs.AppError) {
	var refreshToken string
	var err *errs.AppError
	if refreshToken, err = authtoken.newRefreshToken(); err != nil {
		return "", err
	}
	errr := a.rdb.SAdd(a.ctx, refreshToken, refreshToken, time.Hour).Err()
	if errr != nil {
		return "", errs.NewUnexpectedError("Error generating refresh token" + errr.Error())
	}

	return refreshToken, nil
}

func (a authRepositoryRedisImpl) RefreshTokenExists(refreshToken string) *errs.AppError {
	isExist, err := a.rdb.SIsMember(a.ctx, refreshToken, refreshToken).Result()
	if err != nil {
		return errs.NewUnexpectedError("Error checking if refresh token exists: " + err.Error())
	}
	if !isExist {
		return errs.NewNotFoundError("Refresh token does not exist")
	}
	return nil
}

func NewAuthRepositoryRedisImpl(rdb *redis.Client, ctx context.Context) AuthRepositoryRedis {
	return authRepositoryRedisImpl{rdb, ctx}
}
