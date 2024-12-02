package domain

import (
	"database/sql"
	"errors"
	"github.com/Zzhihon/sso/errs"
	"github.com/jmoiron/sqlx"
	"log"
)

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) FindBy(userID string) (*Login, *errs.AppError) {

	var login Login
	sqlVerify := `SELECT userID, name, role FROM users WHERE userID = ?`

	//通过用户名和密码进行验证是否存在符合的用户
	err := d.client.Get(&login, sqlVerify, userID)
	if err != nil {
		//未找到匹配的用户(id错误)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewUnAuthorizedError("user does not exist")
		} else {
			//其他错误
			log.Println("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewBadGatewayError(err.Error())
		}
	}

	return &login, nil

}

func (d AuthRepositoryDb) GenerateRefreshToken(authtoken AuthToken) (string, *errs.AppError) {
	var err *errs.AppError
	var refreshToken string
	if refreshToken, err = authtoken.newRefreshToken(); err != nil {
		return "", err
	}

	sqlInsert := `INSERT INTO refresh_token (refreshToken) VALUES (?)`
	_, errr := d.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		log.Println("Error while inserting new refresh token from database: " + errr.Error())
		return "", errs.NewBadGatewayError(errr.Error())
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshtoken string) *errs.AppError {
	sqlSelect := `SELECT refreshToken FROM refresh_token WHERE refreshToken = ?`
	var refreshToken string
	err := d.client.Get(&refreshToken, sqlSelect, refreshtoken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Refresh token does not exist")
			return errs.NewNotFoundError("refresh token does not exist")
		} else {
			log.Println("Error while querying refresh token from database: " + err.Error())
			return errs.NewBadGatewayError(err.Error())
		}
	}
	return nil
}

func NewAuthRepositoryDb(client *sqlx.DB) AuthRepository {
	return AuthRepositoryDb{client: client}
}
