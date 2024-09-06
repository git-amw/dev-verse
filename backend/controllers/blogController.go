package controllers

import (
	"net/http"

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
	}
	bc.blogService.CreateBlog(blogDTO)
	ctx.JSON(http.StatusCreated, "Blog is posted successfully")
}

func (bc *blogController) GetAllBlog(ctx *gin.Context) {
	bc.blogService.GetAllBlog()
	// ctx.JSON(http.StatusOK, result)
}
func (bc *blogController) UpdateBlog(ctx *gin.Context) {
	bc.blogService.UpdateBlog()
}
func (bc *blogController) DeleteBlog(ctx *gin.Context) {
	bc.blogService.DeleteBlog()
}

func (bc *blogController) GetAllTags(ctx *gin.Context) {
	result := bc.blogService.GetAllTags()
	ctx.JSON(http.StatusOK, result)
}
