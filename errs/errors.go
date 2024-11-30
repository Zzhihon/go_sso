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
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewUnexpectedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}
