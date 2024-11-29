package service

import (
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/errs"
)

type UserService interface {
	GetAllUsers() ([]domain.User, error)
	GetUser(id string) (*domain.User, *errs.AppError)
}

type DefaultUserService struct {
	repo domain.UserRepository
}

func (s DefaultUserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.FindAll()
}

func (s DefaultUserService) GetUser(id string) (*domain.User, *errs.AppError) {
	return s.repo.ById(id)
}

func NewUserService(repository domain.UserRepository) DefaultUserService {
	return DefaultUserService{repository}
}
