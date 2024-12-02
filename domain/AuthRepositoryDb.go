package domain

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) FindBy(userID string) (*Login, error) {

	var login Login
	sqlVerify := `SELECT userID, name, role FROM users WHERE userID = ?`

	//通过用户名和密码进行验证是否存在符合的用户
	err := d.client.Get(&login, sqlVerify, userID)
	if err != nil {
		//未找到匹配的用户(id错误)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("incorrect user id")
		} else {
			//其他错误
			log.Println("Error while verifying login request from database: " + err.Error())
			return nil, errors.New("unexpected database error")
		}
	}

	return &login, nil

}

func (d AuthRepositoryDb) GenerateRefreshToken(authtoken AuthToken) (string, error) {
	var err error
	var refreshToken string
	if refreshToken, err = authtoken.newRefreshToken(); err != nil {
		return "", err
	}

	sqlInsert := `INSERT INTO refresh_token (refreshToken) VALUES (?)`
	_, err = d.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		log.Println("Error while inserting new refresh token from database: " + err.Error())
		return "", err
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshtoken string) error {
	sqlSelect := `SELECT refreshToken FROM refresh_token WHERE refreshToken = ?`
	var refreshToken string
	err := d.client.Get(&refreshToken, sqlSelect, refreshtoken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Refresh token does not exist")
			return errors.New("invalid refresh token")
		} else {
			log.Println("Error while querying refresh token from database: " + err.Error())
			return errors.New("unexpected database error")
		}
	}
	return nil
}

func NewAuthRepositoryDb(client *sqlx.DB) AuthRepository {
	return AuthRepositoryDb{client: client}
}
