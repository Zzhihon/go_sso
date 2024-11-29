package app

import (
	"github.com/Zhihon/go_sso/domain"
	"github.com/Zhihon/go_sso/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Start() {
	router := mux.NewRouter()

	//初始化一个服务，同时要给这个服务注入依赖(Repo)
	//handler通过Service接口实现业务逻辑，同时依赖Repo来实现与数据库的操作
	ch := UserHandlers{service: service.NewUserService(domain.NewUserRepositoryDb())}

	router.HandleFunc("/Users", ch.getALLUsers).Methods(http.MethodGet)
	//router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	//router.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)
	//router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodPost)
	//router.HandleFunc("/getUser/{username:[0-9]+}", getUser)
	log.Fatal(http.ListenAndServe(":8080", router))
}
