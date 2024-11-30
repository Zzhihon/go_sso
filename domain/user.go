package domain

import (
	"database/sql"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
)

type User struct {
	//sql.NullString用来与mysql的NULL进行映射
	UserID      string         `db:"userID"`
	Password    string         `db:"password"`
	Name        string         `db:"name"`
	Grade       sql.NullString `db:"grade"`
	MajorClass  sql.NullString `db:"majorClass"`
	Email       sql.NullString `db:"email"`
	PhoneNumber sql.NullString `db:"phoneNumber"`
	Status      sql.NullString `db:"status"`
}

type UserRepository interface {
	FindAll(status string) ([]User, *errs.AppError)
	ById(string) (*User, *errs.AppError)
	Update(User, string) (*User, *errs.AppError)
	CheckPassword(User, string) (bool, *errs.AppError)
}

func (u User) StatusAsText() string {
	statusAsText := "active"
	if u.Status.String == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}

func (u User) ToDto() dto.UserResponse {
	return dto.UserResponse{
		UserId:      u.UserID,
		Name:        u.Name,
		Grade:       u.Grade.String,
		MajorClass:  u.MajorClass.String,
		Email:       u.Email.String,
		PhoneNumber: u.PhoneNumber.String,
		Status:      u.StatusAsText(),
	}
}
