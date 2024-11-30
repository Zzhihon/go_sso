package service

import (
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
)

type UserService interface {
	GetAllUsers(status string) ([]dto.UserResponse, *errs.AppError)
	GetUser(id string) (*dto.UserResponse, *errs.AppError)
	UpdateName(r dto.NewUpdateRequest) (*dto.UserResponse, *errs.AppError)
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

func (s DefaultUserService) UpdateName(r dto.NewUpdateRequest) (*dto.UserResponse, *errs.AppError) {
	id := r.UserID
	name := r.Name
	user, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	u := domain.User{
		UserID:      user.UserID,
		Name:        name,
		Grade:       "",
		MajorClass:  "",
		Email:       "",
		PhoneNumber: "",
		Status:      "",
	}

	newUser, err := s.repo.Update(u)
	if err != nil {
		return nil, err
	}
	response := newUser.ToDto()
	return &response, nil
}

func NewUserService(repository domain.UserRepository) DefaultUserService {
	return DefaultUserService{repository}
}
