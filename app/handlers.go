package app

import (
	"encoding/json"
	"github.com/Zhihon/go_sso/service"
	"net/http"
)

type UserHandlers struct {
	service service.UserService
}

func (ch *UserHandlers) getALLUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := ch.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)

}
