package service

import (
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers(status string) ([]dto.UserResponse, *errs.AppError)
	GetUser(id string) (*dto.UserResponse, *errs.AppError)
	Update(r dto.NewUpdateRequest) (*dto.UserResponse, *errs.AppError)
}

type DefaultUserService struct {
	repo      domain.UserRepository
	utilsRepo domain.UtilsRepository
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
	err := checkUserId(id)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.ById(id)

	if err != nil {
		return nil, err
	}
	response := u.ToDto()
	return &response, nil
}

func (s DefaultUserService) Update(r dto.NewUpdateRequest) (*dto.UserResponse, *errs.AppError) {
	id := r.UserID
	err := checkUserId(id)
	if err != nil {
		return nil, err
	}
	var newUser *domain.User
	//寻找有无匹配id的用户
	user, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}

	//Update
	if r.Impl == "Name" {
		if r.Name == "" {
			return nil, errs.NewUnexpectedError("invalid name")
		}
		user.Name = r.Name
	}
	if r.Impl == "Email" {
		if r.Email == "" {
			return nil, errs.NewUnexpectedError("invalid email")
		}
		user.Email.String = r.Email
	}
	if r.Impl == "PhoneNumber" {
		if r.PhoneNumber == "" {
			return nil, errs.NewUnexpectedError("invalid phoneNumber")
		}
		user.PhoneNumber.String = r.PhoneNumber
	}
	if r.Impl == "Role" {
		if r.Role == "" {
			return nil, errs.NewUnexpectedError("invalid role")
		}
		user.Role.String = r.Role
	}
	if r.Impl == "Password" {
		if r.OldPassword == "" {
			return nil, errs.NewUnexpectedError("invalid oldPassword")
		}

		_, ePrr := s.utilsRepo.CheckPassword(id, r.OldPassword)
		if ePrr != nil {
			return nil, ePrr
		}

		if r.NewPassword == "" {
			return nil, errs.NewUnexpectedError("invalid newPassword")
		}

		password, err := hashPassword(r.NewPassword)
		if err != nil {
			return nil, err
		}
		user.Password = password
	}

	newUser, err = s.repo.Update(*user, r.Impl)
	if err != nil {
		return nil, err
	}
	response := newUser.ToDto()
	return &response, nil
}

func hashPassword(password string) (string, *errs.AppError) {
	// bcrypt生成哈希密码，生成的哈希值是一个加盐哈希（salted hash）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.NewUnexpectedError(err.Error())
	}
	return string(hashedPassword), nil
}

func checkUserId(id string) *errs.AppError {
	if id == "" {
		return errs.NewUnexpectedError("invalid user id")
	}
	return nil
}

func NewUserService(repo domain.UserRepository, utilsRepo domain.UtilsRepository) DefaultUserService {
	return DefaultUserService{
		repo:      repo,
		utilsRepo: utilsRepo,
	}
}
