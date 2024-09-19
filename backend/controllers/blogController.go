package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/git-amw/backend/models"
	"github.com/git-amw/backend/services"
)

type BlogController interface {
	CreateBlog(ctx *gin.Context)
	GetAllBlog(ctx *gin.Context)
	UpdateBlog(ctx *gin.Context)
	DeleteBlog(ctx *gin.Context)
	GetAllTags(ctx *gin.Context)
	IncreaseLike(ctx *gin.Context)
}

type blogController struct {
	blogService services.BlogService
}

func NewBlogController(blogService services.BlogService) BlogController {
	return &blogController{
		blogService: blogService,
	}
}

func (bc *blogController) CreateBlog(ctx *gin.Context) {
	var blogDTO models.BlogDTO
	if err := ctx.ShouldBindJSON(&blogDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	success, message := bc.blogService.CreateBlog(blogDTO)
	if success {
		ctx.JSON(http.StatusCreated, gin.H{"message": message})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
	}
}

func (bc *blogController) GetAllBlog(ctx *gin.Context) {
	var result = bc.blogService.GetAllBlog()
	ctx.JSON(http.StatusOK, result)
}
func (bc *blogController) UpdateBlog(ctx *gin.Context) {
	var blogUpdateDTO models.BlogUpdateDTO
	if err := ctx.ShouldBindJSON(&blogUpdateDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	success, message := bc.blogService.UpdateBlog(blogUpdateDTO)
	if success {
		ctx.JSON(http.StatusOK, gin.H{"message": "You've Updated the blog with id : " + strconv.Itoa(blogUpdateDTO.ID)})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"error": message})
}
func (bc *blogController) DeleteBlog(ctx *gin.Context) {
	Id := ctx.Param("id")
	blogId, err := strconv.Atoi(Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error in Converting Id"})
		return
	}
	success, message := bc.blogService.DeleteBlog(blogId)
	if success {
		ctx.JSON(http.StatusOK, gin.H{"message": "You've delted the blog with id : " + Id})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"error": message})
}

func (bc *blogController) GetAllTags(ctx *gin.Context) {
	result := bc.blogService.GetAllTags()
	ctx.JSON(http.StatusOK, result)
}

func (bc *blogController) IncreaseLike(ctx *gin.Context) {
	Id := ctx.Param("id")
	blogId, err := strconv.Atoi(Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error in Converting Id"})
		return
	}
	bc.blogService.IncreaseLike(blogId)
	ctx.JSON(http.StatusOK, gin.H{"message": "You liked the blog with id :" + Id})
}
