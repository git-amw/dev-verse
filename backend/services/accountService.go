package services

import (
	"errors"
	"time"

	"github.com/git-amw/backend/data"
	"github.com/git-amw/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountService interface {
	CreateUser(singnupModel models.SignUp) (bool, string)
	SignInUser(singinModel models.SignIn) (bool, string)
}

type accountService struct {
	DB *gorm.DB
}

func NewAccountService() AccountService {
	return &accountService{
		DB: data.ConnectToDB(),
	}
}

func (as *accountService) CreateUser(signupModel models.SignUp) (bool, string) {
	hashedpassword, err := HashPassword(signupModel.Password)
	if err != nil {
		return false, "Failed to Hash Password"
	}
	signupModel.Password = hashedpassword
	if result := as.DB.Table("sign_ups").Create(&signupModel); result.Error != nil {
		return false, result.Error.Error()
	}
	return true, "User Successfully Created"
}

func (as *accountService) SignInUser(singinModel models.SignIn) (bool, string) {
	var user = struct{ Password string }{}
	if userData := as.DB.Table("sign_ups").Where("email = ?", singinModel.Email).First(&user); userData.Error != nil {
		if errors.Is(userData.Error, gorm.ErrRecordNotFound) {
			return false, "Recode Not Found"
		} else {
			return false, userData.Error.Error()
		}
	} else {
		if !CheckPasswordHash(singinModel.Password, user.Password) {
			return false, "Incorrect Password"
		}
		return GenerateToken(singinModel)
	}

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
