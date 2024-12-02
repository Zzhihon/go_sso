package app

import (
	"encoding/json"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
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
		writeResponse(w, http.StatusBadRequest, errs.NewBadRequestError(err.Error()))
	} else {
		tokens, err := h.service.Login(loginRequest)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, tokens)
		}
	}
}

func (h AuthHandlers) Verify(w http.ResponseWriter, r *http.Request) {
	var verifyRequest dto.VerifyRequest
	var token string
	err := json.NewDecoder(r.Body).Decode(&verifyRequest)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, errs.NewBadRequestError(err.Error()))
	}
	token = verifyRequest.Token
	if token != "" {
		isAuthorized, err := h.service.Verify(token)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, err.AsMessage())
		} else {
			if isAuthorized {
				writeResponse(w, http.StatusOK, "Authorized")
			} else {
				writeResponse(w, http.StatusUnauthorized, err.AsMessage())
			}
		}
	} else {
		writeResponse(w, http.StatusUnauthorized, "Error while verify: "+"missing token")
	}
}

func (h AuthHandlers) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshRequest dto.RefreshRequest
	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, errs.NewBadRequestError(err.Error()))
	} else {
		token, err := h.service.Refresh(refreshRequest)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, token)
		}
	}
}
