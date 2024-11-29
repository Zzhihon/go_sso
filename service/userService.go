package service

import "github.com/Zhihon/go_sso/domain"

type UserService interface {
	GetAllUsers() ([]domain.User, error)
}

type DefaultUserService struct {
	repo domain.UserRepository
}

func (s DefaultUserService) GetAllUsers() ([]domain.User, error) {
	//TODO implement me
	return s.repo.FindAll()
}

func NewUserService(repository domain.UserRepository) DefaultUserService {
	return DefaultUserService{repository}
}
