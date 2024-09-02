package services

import (
	"github.com/git-amw/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type AccountService interface {
	CreateUser(singnupModel models.SignUp) bool
	SignInUser(singinModel models.SignIn) bool
}

type accountService struct{}

func NewAccountService() AccountService {
	return &accountService{}
}

var users = map[string]string{}

func (s *accountService) CreateUser(signupModel models.SignUp) bool {
	hashedpassword, err := HashPassword(signupModel.Password)
	if err != nil {
		return false
	}
	users[signupModel.Email] = hashedpassword
	return true
}

func (s *accountService) SignInUser(singinModel models.SignIn) bool {
	storedPassword, exists := users[singinModel.Email]
	if !exists {
		return false
	}
	if !CheckPasswordHash(singinModel.Password, storedPassword) {
		return false
	}
	return true
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
