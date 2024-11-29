package domain

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type UserRepositoryDb struct {
	client *sql.DB
}

func (d UserRepositoryDb) FindAll() ([]User, error) {
	//sql操作
	findAllSql := "select userID, name from users"
	//返回查询到的rows对象
	rows, err := d.client.Query(findAllSql)
	if err != nil {
		log.Print("Error when query for user table " + err.Error())
		return nil, err
	}
	users := make([]User, 0)
	//查询数据
	for rows.Next() {
		var u User
		//将rows读取到的数据写入User结构体u
		err := rows.Scan(&u.UserID, &u.Name)
		if err != nil {
			log.Print("Error when scanning user table " + err.Error())
			return nil, err
		}
		users = append(users, u)
	}
	//此时的数据库只初始化了name和userID的字段，其他字段还没涉及到sql查询
	//所以这里返回的结构体会包含null值
	return users, nil
}

func NewUserRepositoryDb() UserRepositoryDb {
	//远程连接到数据库
	client, err := sql.Open("mysql", "root:7tvkPQzKGe1Syv5E@tcp(127.0.0.1:3306)/sso")
	if err != nil {
		panic(err)
	}
	//配置mysql连接池
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return UserRepositoryDb{client}
}
