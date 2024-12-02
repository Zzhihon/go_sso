package errs

import "net/http"

type AppError struct {
	Code    int    `json:",omitempty"`
	Message string `json:"message"`
}

// helper function
// 用作输出AppError存异常信息
func (err AppError) AsMessage() *AppError {
	return &AppError{
		Message: err.Message}
}

// 自定义异常处理

// 404 sql查询失败
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

// 500
func NewUnexpectedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

// api路径出错
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

// 401 用户名或密码出错
func NewUnAuthorizedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

// 502 sql查询失败
func NewBadGatewayError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusBadGateway,
	}
}
