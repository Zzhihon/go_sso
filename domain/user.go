package domain

type User struct {
	UserID string
	Name   string
}

type UserRepository interface {
	FindAll() ([]User, error)
}
