package app

import (
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"
)

func Start() {
	router := mux.NewRouter()

	//初始化一个服务，同时要给这个服务注入依赖(Repo)
	//handler通过Service接口实现业务逻辑，同时依赖Repo来实现与数据库的操作
	ch := UserHandlers{service: service.NewUserService(domain.NewUserRepositoryDb(getDBClient()))}
	ah := AuthHandlers{service: service.NewAuthService(domain.NewAuthRepositoryDb(getDBClient()))}

	router.HandleFunc("/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/verify", ah.Verify).Methods(http.MethodPost)

	router.HandleFunc("/Update/{impl:[a-zA-Z0-9]+}", ch.update).Methods((http.MethodPost))
	router.HandleFunc("/GetUser/{user_id:[0-9]+}", ch.getUser).Methods("GET")
	router.HandleFunc("/Users", ch.getALLUsers).Methods(http.MethodGet)
	//router.HandleFunc("/getUser/{username:[0-9]+}", getUser)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getDBClient() *sqlx.DB {
	//远程连接到数据库
	client, err := sqlx.Open("mysql", "root:7tvkPQzKGe1Syv5E@tcp(127.0.0.1:3306)/sso")
	if err != nil {
		panic(err)
	}
	//配置mysql连接池
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
