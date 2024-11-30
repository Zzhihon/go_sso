package app

import (
	"encoding/json"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/service"
	"net/http"
)

type AuthHandlers struct {
	service service.AuthService
}

func (h AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "Error while decoding: "+err.Error())
	} else {
		token, err := h.service.Login(loginRequest)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, "Error while login: "+err.Error())
		} else {
			writeResponse(w, http.StatusOK, token)
		}
	}
}
