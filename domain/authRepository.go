package domain

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type AuthRepository interface {
	FindBy(userID string, password string) (*Login, error)
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) FindBy(userID string, password string) (*Login, error) {
	var login Login
	sqlVerify := `SELECT userID, password FROM users WHERE userID = ? and password = ?`

	//通过用户名和密码进行验证是否存在符合的用户
	//后续加上哈希加密验证
	err := d.client.Get(&login, sqlVerify, userID, password)
	if err != nil {
		//未找到匹配的用户(密码或id错误)
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		} else {
			//其他错误
			log.Println("Error while verifying login request from database: " + err.Error())
			return nil, errors.New("Unexpected database error")
		}
	}

	return &login, nil

}

func NewAuthRepositoryDb(client *sqlx.DB) AuthRepository {
	return AuthRepositoryDb{client: client}
}
