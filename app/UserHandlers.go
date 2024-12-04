package app

import (
	"encoding/json"
	"encoding/xml"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/service"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandlers struct {
	service service.UserService
}

func (ch *UserHandlers) getALLUsers(w http.ResponseWriter, r *http.Request) {

	status := r.URL.Query().Get("status")
	users, _ := ch.service.GetAllUsers(status)

	if r.Header.Get("Content-Type") == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(users)
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func (ch *UserHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["user_id"]

	user, err := ch.service.GetUser(id)

	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, user)
	}
}

func (ch *UserHandlers) update(w http.ResponseWriter, r *http.Request) {
	var request dto.NewUpdateRequest
	//处理路径中要执行的操作
	vars := mux.Vars(r)
	impl := vars["impl"]

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		request.Impl = impl
		user, err := ch.service.Update(request)
		if err != nil {
			writeResponse(w, err.Code, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, user)
		}
	}
}

func (ch *UserHandlers) IsEmailValid(w http.ResponseWriter, r *http.Request) {
	var request dto.CheckEmailRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		token, err := ch.service.IsEmailValid(request)
		if err != nil {
			writeResponse(w, err.Code, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, token)
		}
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
