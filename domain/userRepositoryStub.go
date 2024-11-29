package domain

type UserRepositoryStub struct {
	users []User
}

func (s UserRepositoryStub) FindAll() ([]User, error) {
	return s.users, nil
}

func NewUserRepositoryStub() UserRepositoryStub {
	users := []User{
		{"20233802086", "hzh"},
		{"20233808888", "qwq"},
	}

	return UserRepositoryStub{users}
}
