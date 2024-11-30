package domain

import (
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
)

type User struct {
	UserID      string `db:"userID"`
	Name        string `db:"name"`
	Grade       string `db:"grade"`
	MajorClass  string `db:"majorClass"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phoneNumber"`
	Status      string `db:"status"`
}

func (u User) StatusAsText() string {
	statusAsText := "active"
	if u.Status == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}

func (u User) ToDto() dto.UserResponse {
	return dto.UserResponse{
		UserId:      u.UserID,
		Name:        u.Name,
		Grade:       u.Grade,
		MajorClass:  u.MajorClass,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Status:      u.StatusAsText(),
	}
}

type UserRepository interface {
	FindAll(status string) ([]User, *errs.AppError)
	ById(string) (*User, *errs.AppError)
	Update(User) (*User, *errs.AppError)
}
