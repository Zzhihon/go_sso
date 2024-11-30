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
	Usersql := "select userID, name from users where userID = ?"

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

func (d UserRepositoryDb) Update(u User) (*User, *errs.AppError) {
	// 更新 email 的 SQL 查询
	query := "UPDATE users SET name = ? WHERE userID = ?;"

	// 使用 Exec 执行更新操作
	result, err := d.client.Exec(query, u.Name, u.UserID)
	if err != nil {
		log.Fatal(err)
	}

	// 获取更新的行数
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if affectedRows == 0 {
		return nil, errs.NewUnexpectedError("No rows were updated")
	}

	return &u, nil
}

func NewUserRepositoryDb(client *sqlx.DB) UserRepositoryDb {
	return UserRepositoryDb{client: client}
}
