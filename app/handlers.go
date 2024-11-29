package app

import (
	"encoding/json"
	"encoding/xml"
	"github.com/Zhihon/go_sso/service"
	"net/http"
)

type UserHandlers struct {
	service service.UserService
}

func (ch *UserHandlers) getALLUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := ch.service.GetAllUsers()
	if r.Header.Get("Content-Type") == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(users)
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
