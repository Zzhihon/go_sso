package domain

import "database/sql"

type User struct {
	UserID      string         `db: "userID"`
	Name        string         `db: "name"`
	Grade       sql.NullString `db: "grade"`
	MajorClass  sql.NullString `db: "majorClass"`
	Email       sql.NullString `db: "email"`
	PhoneNumber sql.NullString `db: "phoneNumber"`
}

type UserRepository interface {
	FindAll() ([]User, error)
}
