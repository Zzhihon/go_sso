package service

import (
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
)

type UserService interface {
	GetAllUsers(status string) ([]dto.UserResponse, *errs.AppError)
	GetUser(id string) (*dto.UserResponse, *errs.AppError)
	Update(r dto.NewUpdateRequest) (*dto.UserResponse, *errs.AppError)
}

type DefaultUserService struct {
	repo domain.UserRepository
}

func (s DefaultUserService) GetAllUsers(status string) ([]dto.UserResponse, *errs.AppError) {
	if status == "active" {
		status = "1"
	} else if status == "inactive" {
		status = "0"
	} else {
		status = ""
	}

	user, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}
	response := make([]dto.UserResponse, 0)
	for _, u := range user {
		response = append(response, u.ToDto())
	}

	return response, err
}

func (s DefaultUserService) GetUser(id string) (*dto.UserResponse, *errs.AppError) {
	u, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	response := u.ToDto()
	return &response, nil
}

func (s DefaultUserService) Update(r dto.NewUpdateRequest) (*dto.UserResponse, *errs.AppError) {
	id := r.UserID
	//寻找有无匹配id的用户
	user, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}

	//Update
	if r.Impl == "Name" {
		user.Name = r.Name
	}
	if r.Impl == "Email" {
		user.Email.String = r.Email
	}
	if r.Impl == "PhoneNumber" {
		user.PhoneNumber.String = r.PhoneNumber
	}

	newUser, err := s.repo.Update(*user, r.Impl)
	if err != nil {
		return nil, err
	}
	response := newUser.ToDto()
	return &response, nil
}

func NewUserService(repository domain.UserRepository) DefaultUserService {
	return DefaultUserService{repository}
}
