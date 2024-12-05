package app

import (
	"context"
	"fmt"
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/logger"
	"github.com/Zzhihon/sso/service"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

func Start() {
	sanityCheck()
	router := mux.NewRouter()

	//初始化一个服务，同时要给这个服务注入依赖(Repo)
	//handler通过Service接口实现业务逻辑，同时依赖Repo来实现与数据库的操作
	ch := UserHandlers{service: service.NewUserService(domain.NewUserRepositoryDb(getDBClient()), domain.NewUtilsRepositoryDb(getDBClient()), domain.NewRedisRepositoryImpl(initRedis(), context.Background()))}
	ah := AuthHandlers{service: service.NewAuthService(domain.NewAuthRepositoryDb(getDBClient()), domain.NewUtilsRepositoryDb(getDBClient()), domain.NewRedisRepositoryImpl(initRedis(), context.Background()))}

	router.HandleFunc("/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/verify", ah.Verify).Methods(http.MethodPost)
	router.HandleFunc("/refresh", ah.Refresh).Methods(http.MethodPost)

	router.HandleFunc("/heartbeat", ch.HeartBeat).Methods(http.MethodPost)
	router.HandleFunc("/Code", ch.IsEmailValid).Methods(http.MethodPost)
	router.HandleFunc("/Update/{impl:[a-zA-Z0-9]+}", ch.update).Methods(http.MethodPost)
	router.HandleFunc("/GetUser/{user_id:[0-9]+}", ch.getUser).Methods(http.MethodGet)
	router.HandleFunc("/Users", ch.getALLUsers).Methods(http.MethodGet)
	//router.HandleFunc("/getUser/{username:[0-9]+}", getUser)
	SERVER_PORT := os.Getenv("DB_USER")

	log.Fatal(http.ListenAndServe(SERVER_PORT, router))
}

func getDBClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// postgresql数据库连接信息
	//dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",dbAddr, dbPort, dbUser, dbPasswd, dbName)
	//db, err := sqlx.Connect("postgres", dsn) // 使用 PostgreSQL 驱动
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}

func initRedis() *redis.Client {
	var rdb *redis.Client
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func sanityCheck() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger.Info("load .env success")
	envProps := []string{
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Error(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}
