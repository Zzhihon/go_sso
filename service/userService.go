package service

import (
	"bytes"
	"github.com/Zzhihon/sso/domain"
	"github.com/Zzhihon/sso/dto"
	"github.com/Zzhihon/sso/errs"
	"github.com/Zzhihon/sso/utils"
	"github.com/go-gomail/gomail"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"math/rand"
	"time"
)

type UserService interface {
	GetAllUsers(status string) ([]dto.UserResponse, *errs.AppError)
	GetUser(id string) (*dto.UserResponse, *errs.AppError)
	Update(r dto.NewUpdateRequest) (*dto.UserResponse, *errs.AppError)
	IsEmailValid(r dto.CheckEmailRequest) (string, *errs.AppError)
}

type DefaultUserService struct {
	repo      domain.UserRepository
	utilsRepo domain.UtilsRepository
	redis     domain.RedisRepository
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

func (s DefaultUserService) IsEmailValid(r dto.CheckEmailRequest) (string, *errs.AppError) {
	var err *errs.AppError

	err = checkUserId(r.UserID)
	if err != nil {
		return "", err
	}
	
	err = s.repo.IsEmailValid(r.UserID, r.Email)
	if err != nil {
		return "", err
	} else {
		//生成随机码，发送邮箱，并让code验证码存储在redis里
		var code string
		code = generateRandomString()
		errr := s.redis.StoreUserCode(r.UserID, code)
		if errr != nil {
			return "", errr
		}

		var u *domain.User
		u, err = s.repo.ById(r.UserID)
		if err != nil {
			return "", err
		}
		err = sendEmail(*u, code)
		if err != nil {
			return "", err
		}
		return code, nil
	}
}

func sendEmail(u domain.User, code string) *errs.AppError {

	// 1. 读取并解析外部 HTML 模板文件
	tmpl, err := template.ParseFiles("templates/email_template.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return errs.NewUnexpectedError("Error parsing template" + err.Error())
	}

	// 2. 渲染模板到缓冲区
	var body bytes.Buffer
	err = tmpl.Execute(&body, dto.EmailData{
		Name: u.Name,
		Code: code,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
		return errs.NewUnexpectedError("Error executing template" + err.Error())
	}

	//3. 使用gmail发送邮件
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@esaps.net") // 发件人
	m.SetHeader("To", u.Email.String)         // 收件人
	m.SetHeader("Subject", "Test Email")      // 主题
	m.SetBody("text/html", body.String())     //正文

	// SMTP 服务器配置
	d := gomail.NewDialer("smtp.esaps.net", 587, "no-reply@esaps.net", utils.Email_Password)
	errr := d.DialAndSend(m)
	if errr != nil {
		return errs.NewUnexpectedError("Error while sending email" + errr.Error())
	}

	// 发送邮件
	return nil
}

func generateRandomString() string {
	const charset = utils.Charset
	const length = 6
	rand.Seed(time.Now().UnixNano()) // 使用当前时间的纳秒数作为种子

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
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
		//if r.OldPassword == "" {
		//	return nil, errs.NewUnexpectedError("invalid oldPassword")
		//}
		//
		//_, ePrr := s.utilsRepo.CheckPassword(id, r.OldPassword)
		//if ePrr != nil {
		//	return nil, ePrr
		//}

		err := s.redis.IsCodeExists(id, r.Code)
		if err != nil {
			return nil, err
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

func NewUserService(repo domain.UserRepository, utilsRepo domain.UtilsRepository, redis domain.RedisRepository) DefaultUserService {
	return DefaultUserService{
		repo:      repo,
		utilsRepo: utilsRepo,
		redis:     redis,
	}
}
