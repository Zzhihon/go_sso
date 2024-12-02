package domain

import (
	"database/sql"
	"errors"
	"github.com/Zzhihon/sso/errs"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UtilsRepository interface {
	CheckPassword(id string, inputpassword string) (bool, *errs.AppError)
}

type UtilsRepositoryDb struct {
	client *sqlx.DB
}

func (d UtilsRepositoryDb) CheckPassword(id string, inputpassword string) (bool, *errs.AppError) {
	Usersql := "select password from users where userID = ?"
	var password string
	err := d.client.Get(&password, Usersql, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.NewUnAuthorizedError("user id does not exit")
		} else {
			log.Println("Error while checking password from database: " + err.Error())
			return false, errs.NewBadGatewayError(err.Error())
		}
	}

	pErr := bcrypt.CompareHashAndPassword([]byte(password), []byte(inputpassword))
	if pErr != nil {
		return false, errs.NewUnAuthorizedError("password does not match")
	}
	return true, nil
}

func NewUtilsRepositoryDb(client *sqlx.DB) UtilsRepositoryDb {
	return UtilsRepositoryDb{client: client}
}
