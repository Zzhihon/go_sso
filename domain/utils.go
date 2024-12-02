package domain

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UtilsRepository interface {
	CheckPassword(id string, inputpassword string) (bool, error)
}

type UtilsRepositoryDb struct {
	client *sqlx.DB
}

func (d UtilsRepositoryDb) CheckPassword(id string, inputpassword string) (bool, error) {
	Usersql := "select password from users where userID = ?"
	var password string
	err := d.client.Get(&password, Usersql, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("userid not exist")
		} else {
			log.Println("Error while checking password from database: " + err.Error())
			return false, errors.New("unexpected database error")
		}
	}

	pErr := bcrypt.CompareHashAndPassword([]byte(password), []byte(inputpassword))
	if pErr != nil {
		return false, errors.New("password not correct")
	}
	return true, nil
}

func NewUtilsRepositoryDb(client *sqlx.DB) UtilsRepositoryDb {
	return UtilsRepositoryDb{client: client}
}
