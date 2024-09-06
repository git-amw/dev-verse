package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/git-amw/backend/models"
	"github.com/git-amw/backend/services"
	"net/http"
)

type AccountController interface {
	CreateUser(ctx *gin.Context)
	SignInUser(ctx *gin.Context)
}

type accountController struct {
	accountService services.AccountService
}

func NewAccountController(accountService services.AccountService) AccountController {
	return &accountController{
		accountService: accountService,
	}
}

func (ac *accountController) CreateUser(ctx *gin.Context) {
	var singupModel models.SignUp
	if err := ctx.ShouldBindBodyWithJSON(&singupModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, message := ac.accountService.CreateUser(singupModel)
	if ok {
		ctx.JSON(http.StatusCreated, gin.H{"message": message})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
	}
}

func (ac *accountController) SignInUser(ctx *gin.Context) {
	var singinModel models.SignIn
	if err := ctx.ShouldBindJSON(&singinModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, message := ac.accountService.SignInUser(singinModel)
	if ok {
		ctx.JSON(http.StatusOK, gin.H{"token": message})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": message})
	}
}
