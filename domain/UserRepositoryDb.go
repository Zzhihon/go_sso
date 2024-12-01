package domain

import (
	"database/sql"
	"fmt"
	"github.com/Zzhihon/sso/errs"
	"github.com/Zzhihon/sso/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepositoryDb struct {
	client *sqlx.DB
}

func (d UserRepositoryDb) FindAll(status string) ([]User, *errs.AppError) {
	var err error
	users := make([]User, 0)

	if status == "" {
		findAllSql := "select userID, name from users"
		err = d.client.Select(&users, findAllSql)
	} else {
		//筛选出status为某一特定状态的所有用户
		findAllSql := "select userID, name from users where status = ?"
		err = d.client.Select(&users, findAllSql, status)
	}
	if err != nil {
		logger.Error("Error while querying user table " + err.Error())
		return nil, errs.NewNotFoundError("Unexpected database error")
	}
	//此时的数据库只初始化了name和userID的字段，其他字段还没涉及到sql查询
	//所以这里返回的结构体会包含null值
	return users, nil
}

func (d UserRepositoryDb) ById(id string) (*User, *errs.AppError) {
	Usersql := "select userID, name, email, phoneNumber, grade, majorClass, status from users where userID = ?"

	var u User
	err := d.client.Get(&u, Usersql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("User not found")
		} else {
			logger.Error("Error while querying user table " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	return &u, nil
}

func (d UserRepositoryDb) Update(u User, imple string) (*User, *errs.AppError) {
	var query string
	var s string
	//用imple识别用户要修改的字段
	if imple == "Name" {
		query = "UPDATE users SET name = ? WHERE userID = ?;"
		s = u.Name
	}
	if imple == "Email" {
		query = "UPDATE users SET email = ? WHERE userID = ?;"
		s = u.Email.String
	}
	if imple == "PhoneNumber" {
		query = "UPDATE users SET phoneNumber = ? WHERE userID = ?;"
		s = u.PhoneNumber.String
	}
	if imple == "Password" {
		query = "UPDATE users SET password = ? WHERE userID = ?;"
		s = u.Password
	}
	if imple == "Role" {
		query = "UPDATE users SET role = ? WHERE userID = ?;"
		s = u.Role.String
	}

	// 使用 Exec 执行更新操作
	//UserID用来锁定row
	//结构体的字段会和数据库的字段进行映射 确保后端和数据都会同步更新
	result, err := d.client.Exec(query, s, u.UserID)
	if err != nil {
		log.Fatal(err)
	}

	// 获取更新的行数
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	//异常处理
	if affectedRows == 0 {
		return nil, errs.NewUnexpectedError("No rows were updated")
	}

	return &u, nil
}

func (d UserRepositoryDb) CheckPassword(u User, originPassword string) (bool, *errs.AppError) {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(originPassword))
	if err != nil {
		fmt.Println("Password does not match")
		return false, errs.NewUnexpectedError("Password does not match")
	}

	return true, nil

}

func NewUserRepositoryDb(client *sqlx.DB) UserRepositoryDb {
	return UserRepositoryDb{client: client}
}
