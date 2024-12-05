package app

import (
	"encoding/json"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type UserHandlers struct {
	service service.UserService
}

func (ch *UserHandlers) getALLUsers(w http.ResponseWriter, r *http.Request) {
	// 获取分页参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	status := r.URL.Query().Get("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 0 {
		pageSize = 50
	}

	request := dto.GetAllUsers{
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}

	users, errr := ch.service.GetAllUsers(request)
	if errr != nil {
		log.Println(errr)
		writeResponse(w, errr.Code, errr.AsMessage())
	} else {
		writeResponse(w, http.StatusOK, users)
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

func (ch *UserHandlers) HeartBeat(w http.ResponseWriter, r *http.Request) {
	var request dto.OnlineUsers

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		var res *dto.UserStateResponse
		res, err := ch.service.UserOnline(request.UserID)
		if err != nil {
			writeResponse(w, err.Code, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, res)
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
