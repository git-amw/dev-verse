package services

import (
	"time"

	"github.com/git-amw/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AccountService interface {
	CreateUser(singnupModel models.SignUp) bool
	SignInUser(singinModel models.SignIn) (bool, string)
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

func (s *accountService) SignInUser(singinModel models.SignIn) (bool, string) {
	storedPassword, exists := users[singinModel.Email]
	if !exists {
		return false, "User Already Exists"
	}
	if !CheckPasswordHash(singinModel.Password, storedPassword) {
		return false, "Incorrect Password"
	}
	return GenerateToken(singinModel)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(singinModel models.SignIn) (bool, string) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": singinModel.Email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	})
	token, err := generateToken.SignedString([]byte("SECRET_KEY"))
	if err != nil {
		return false, "Failed to generate token"
	}
	return true, token
}
