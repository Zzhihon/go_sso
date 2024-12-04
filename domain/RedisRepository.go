package domain

import (
	"context"
	"github.com/Zzhihon/sso/errs"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepository interface {
	GenerateRefreshToken(token AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
	StoreUserCode(userID string, token string) *errs.AppError
	IsCodeExists(userID string, token string) *errs.AppError
}

type RedisRepositoryImpl struct {
	rdb *redis.Client
	ctx context.Context
}

func (a RedisRepositoryImpl) GenerateRefreshToken(authtoken AuthToken) (string, *errs.AppError) {
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

func (a RedisRepositoryImpl) RefreshTokenExists(refreshToken string) *errs.AppError {
	isExist, err := a.rdb.SIsMember(a.ctx, refreshToken, refreshToken).Result()
	if err != nil {
		return errs.NewUnexpectedError("Error checking if refresh token exists: " + err.Error())
	}
	if !isExist {
		return errs.NewNotFoundError("Refresh token does not exist")
	}
	return nil
}

func (a RedisRepositoryImpl) StoreUserCode(userID string, code string) *errs.AppError {
	err := a.rdb.Set(a.ctx, userID, code, time.Minute*10).Err()
	if err != nil {
		return errs.NewUnexpectedError("Error while store set" + err.Error())
	} else {
		return nil
	}
}

func (a RedisRepositoryImpl) IsCodeExists(userID string, code string) *errs.AppError {
	// 查询 refresh_token
	code, err := a.rdb.Get(a.ctx, userID).Result()
	if err == redis.Nil {
		return errs.NewNotFoundError("Token does not exist " + err.Error())
	}
	if err != nil {
		return errs.NewUnexpectedError("Error while checking token " + err.Error())
	}
	return nil
}

func NewRedisRepositoryImpl(rdb *redis.Client, ctx context.Context) RedisRepository {
	return RedisRepositoryImpl{rdb, ctx}
}
