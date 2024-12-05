package domain

import (
	"database/sql"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
)

type User struct {
	//sql.NullString用来与mysql的NULL进行映射
	UserID      string         `db:"username"`
	Password    string         `db:"password"`
	Name        string         `db:"name"`
	Grade       sql.NullString `db:"grade"`
	MajorClass  sql.NullString `db:"major_class"`
	Email       sql.NullString `db:"email"`
	PhoneNumber sql.NullString `db:"phone_number"`
	//Status      sql.NullString `db:"status"`
	Role        sql.NullString `db:"role"`
	IsActive    bool           `db:"is_active"`
	IsSuperuser bool           `db:"is_superuser"`
	IsStaff     bool           `db:"is_staff"`
}

type UserRepository interface {
	FindAll(status string, page int, pageSize int) ([]User, *errs.AppError)
	ById(string) (*User, *errs.AppError)
	Update(User, string) (*User, *errs.AppError)
	IsEmailValid(string, email string) *errs.AppError
}

//func (u User) StatusAsText() string {
//	statusAsText := "active"
//	if u.Status.String == "0" {
//		statusAsText = "inactive"
//	}
//	return statusAsText
//}

func (u User) ToDto() dto.UserResponse {
	return dto.UserResponse{
		UserId:      u.UserID,
		Name:        u.Name,
		Grade:       u.Grade.String,
		MajorClass:  u.MajorClass.String,
		Email:       u.Email.String,
		PhoneNumber: u.PhoneNumber.String,
		IsActive:    u.IsActive,
		IsSuperuser: u.IsSuperuser,
		IsStaff:     u.IsStaff,
		//Status:      u.StatusAsText(),
	}
}
