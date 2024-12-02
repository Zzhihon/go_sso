package domain

import (
	"database/sql"
	"github.com/Zzhihon/sso/errs"
	"github.com/Zzhihon/sso/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type UserRepositoryDb struct {
	client *sqlx.DB
}

func (d UserRepositoryDb) FindAll(status string) ([]User, *errs.AppError) {
	var err error
	users := make([]User, 0)

	if status == "" {
		findAllSql := "select username, name from account_customuser"
		err = d.client.Select(&users, findAllSql)

	} else {
		//筛选出status为某一特定状态的所有用户
		findAllSql := "select username, name from account_customuser where status = ?"
		err = d.client.Select(&users, findAllSql, status)
	}
	if err != nil {
		logger.Error("Error while querying user table " + err.Error())
		return nil, errs.NewNotFoundError(err.Error())
	}
	//此时的数据库只初始化了name和userID的字段，其他字段还没涉及到sql查询
	//所以这里返回的结构体会包含null值
	return users, nil
}

func (d UserRepositoryDb) ById(id string) (*User, *errs.AppError) {
	Usersql := "select username, name, email, phone_number, grade, major_class, status from account_customuser where username = ?"

	var u User
	err := d.client.Get(&u, Usersql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("user id does not exist")
		} else {
			logger.Error("Error while querying user table " + err.Error())
			return nil, errs.NewBadGatewayError(err.Error())
		}
	}

	return &u, nil
}

func (d UserRepositoryDb) Update(u User, imple string) (*User, *errs.AppError) {
	var query string
	var s string
	//用imple识别用户要修改的字段
	if imple == "Name" {
		query = "UPDATE account_customuser SET name = ? WHERE username = ?;"
		s = u.Name
	}
	if imple == "Email" {
		query = "UPDATE account_customuser SET email = ? WHERE username = ?;"
		s = u.Email.String
	}
	if imple == "PhoneNumber" {
		query = "UPDATE account_customuser SET phone_number = ? WHERE username = ?;"
		s = u.PhoneNumber.String
	}
	if imple == "Password" {
		query = "UPDATE account_customuser SET password = ? WHERE username = ?;"
		s = u.Password
	}
	if imple == "Role" {
		query = "UPDATE account_customuser SET role = ? WHERE username = ?;"
		s = u.Role.String
	}

	// 使用 Exec 执行更新操作
	//UserID用来锁定row
	//结构体的字段会和数据库的字段进行映射 确保后端和数据都会同步更新
	result, err := d.client.Exec(query, s, u.UserID)
	if err != nil {
		log.Fatal(err)
		return nil, errs.NewUnexpectedError(err.Error())
	}

	// 获取更新的行数
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	//如果没有更新的话
	if affectedRows == 0 {
		return &u, nil
	}

	return &u, nil
}

func NewUserRepositoryDb(client *sqlx.DB) UserRepositoryDb {
	return UserRepositoryDb{client: client}
}
