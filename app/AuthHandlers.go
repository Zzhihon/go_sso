package app

import (
	"encoding/json"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
	"github.com/Zzhihon/sso/service"
	"github.com/Zzhihon/sso/utils"
	"log"
	"net/http"
	"time"
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
		log.Println("success login service")
		log.Println(tokens)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, err.AsMessage())
		} else {
			// 设置 access_token 到 cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    tokens.AccessToken,
				HttpOnly: true,  // 安全标记，防止通过 JavaScript 访问
				Secure:   false, // 在开发时可以设置为 false，生产环境使用 true
				Path:     "/",
				Expires:  time.Now().Add(utils.ACCESS_TOKEN_DURATION), // 设置有效期
			})

			// 设置 refresh_token 到 cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    tokens.RefreshToken,
				HttpOnly: true,
				Secure:   false,
				Path:     "/",
				Expires:  time.Now().Add(utils.REFRESH_TOKEN_DURATION), // refresh token 可以设置较长有效期
			})

			// 发送响应
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
