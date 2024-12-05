package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zzhihon/sso/errs"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisRepository interface {
	GenerateRefreshToken(token AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
	StoreUserCode(userID string, token string) *errs.AppError
	IsCodeExists(userID string, token string) *errs.AppError
	StoreUserOnline(userID string) *errs.AppError
	GetOnlineUsers() (int, *errs.AppError)
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
	errr := a.rdb.Set(a.ctx, refreshToken, "refreshToken", time.Minute).Err()
	if errr != nil {
		return "", errs.NewUnexpectedError("Error generating refresh token" + errr.Error())
	}

	return refreshToken, nil
}

func (a RedisRepositoryImpl) RefreshTokenExists(refreshToken string) *errs.AppError {
	isExist, err := a.rdb.Exists(a.ctx, refreshToken).Result()
	if err != nil {
		return errs.NewUnexpectedError("Error checking if refresh token exists: " + err.Error())
	}
	if isExist < 1 {
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

func (a RedisRepositoryImpl) StoreUserOnline(userID string) *errs.AppError {
	//redis设置60s过期
	//如果键已经存在，覆盖原来的值，不用返回报错，
	//若键不存在，则存到redis里面
	err := a.rdb.Set(a.ctx, "**"+userID, "online", time.Minute).Err()
	if err != nil {
		return errs.NewUnexpectedError("Error while store set" + err.Error())
	} else {
		return nil
	}
}

func (a RedisRepositoryImpl) GetOnlineUsers() (int, *errs.AppError) {
	var keys []string
	cursor := uint64(0)
	var pattern string
	pattern = "\\*\\**"

	// 使用 SCAN 命令遍历 Redis 键
	for {
		// SCAN 命令返回游标和匹配的键
		result, newCursor, err := a.rdb.Scan(a.ctx, cursor, pattern, 50).Result()
		if err != nil {
			return -1, errs.NewUnexpectedError("Error while fetching online users" + err.Error())
		}

		// 将结果追加到 keys 切片中
		keys = append(keys, result...)

		// 如果游标为 0，表示扫描完成
		if newCursor == 0 {
			break
		}

		// 更新游标，继续扫描
		cursor = newCursor
	}

	return len(keys), nil
}

func (a RedisRepositoryImpl) IsCodeExists(userID string, code string) *errs.AppError {
	storedToken, err := a.getToken(userID)
	if err != nil {
		log.Fatalf("Error checking token: %v", err)
	}

	// 验证 token 是否匹配
	if validateToken(code, storedToken) {
		fmt.Println("Token is valid and matches.")
		return nil
	} else {
		fmt.Println("Token is invalid or does not match.")
		return errs.NewUnexpectedError("Token is invalid or does not match.")
	}
}

func (a RedisRepositoryImpl) getToken(id string) (string, *errs.AppError) {
	// 使用 GET 命令获取 id 对应的 token
	code, err := a.rdb.Get(a.ctx, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 如果键不存在，返回空字符串
			return "", nil
		}
		return "", errs.NewUnexpectedError("Error while fetching token " + err.Error())
	}
	return code, nil
}

func validateToken(providedToken, storedToken string) bool {
	// 比较提供的 token 和存储的 token 是否一致
	return providedToken == storedToken
}

func NewRedisRepositoryImpl(rdb *redis.Client, ctx context.Context) RedisRepository {
	return RedisRepositoryImpl{rdb, ctx}
}
