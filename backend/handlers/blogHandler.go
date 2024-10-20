package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/git-amw/backend/models"
	"github.com/git-amw/backend/services"
)

type BlogHandlerProvider interface {
	CreateBlog(ctx *gin.Context)
	GetAllBlog(ctx *gin.Context)
	UpdateBlog(ctx *gin.Context)
	DeleteBlog(ctx *gin.Context)
	DeleteTagFromBlog(ctx *gin.Context)
	GetAllTags(ctx *gin.Context)
	IncreaseLike(ctx *gin.Context)
	SearchTags(ctx *gin.Context)
	SearchBlog(ctx *gin.Context)
}

type BlogHandler struct {
	blogService services.BlogServiceProvider
}

func NewBlogHandler(blogService services.BlogServiceProvider) BlogHandlerProvider {
	return &BlogHandler{
		blogService: blogService,
	}
}

func (bh *BlogHandler) CreateBlog(ctx *gin.Context) {
	userId, ok := ctx.Get("userid")
	if !ok {
		log.Println("User Id not found -blog Handler")
		return
	}
	var blogData models.Blog
	if err := ctx.ShouldBindJSON(&blogData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := ConvertId(userId)
	success, message := bh.blogService.CreateBlog(blogData, uid)
	if success {
		ctx.JSON(http.StatusCreated, gin.H{"message": message})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
	}

}

func (bh *BlogHandler) GetAllBlog(ctx *gin.Context) {
	var result = bh.blogService.GetAllBlog()
	ctx.JSON(http.StatusOK, result)
}

func (bh *BlogHandler) UpdateBlog(ctx *gin.Context) {
	var updateData models.BlogUpdate
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	success, message := bh.blogService.UpdateBlog(updateData)
	if success {
		ctx.JSON(http.StatusOK, gin.H{"message": "You've Updated the blog with id : " + strconv.Itoa(int(updateData.ID))})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"error": message})
}

func (bh *BlogHandler) DeleteBlog(ctx *gin.Context) {
	Id := ctx.Param("id")
	blogId, err := strconv.Atoi(Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error in Converting Id"})
		return
	}
	success, message := bh.blogService.DeleteBlog(blogId)
	if success {
		ctx.JSON(http.StatusOK, gin.H{"message": "You've delted the blog with id : " + Id})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"error": message})
}

func (bh *BlogHandler) DeleteTagFromBlog(ctx *gin.Context) {

}

func (bh *BlogHandler) GetAllTags(ctx *gin.Context) {
	result := bh.blogService.GetAllTags()
	ctx.JSON(http.StatusOK, result)
}

func (bh *BlogHandler) IncreaseLike(ctx *gin.Context) {
	Id := ctx.Param("id")
	blogId, err := strconv.Atoi(Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error in Converting Id"})
		return
	}
	bh.blogService.IncreaseLike(blogId)
	ctx.JSON(http.StatusOK, gin.H{"message": "You liked the blog with id :" + Id})
}

func (bh *BlogHandler) SearchTags(ctx *gin.Context) {
	tagValue := ctx.Query("tagValue")
	log.Println(tagValue)
	if tagValue == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}
	tagSearchResponse := bh.blogService.SearchTags(tagValue)
	if tagSearchResponse.TagId == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No tag found"})
		return
	}
	ctx.JSON(http.StatusFound, gin.H{"tag": tagSearchResponse})
}

func (bh *BlogHandler) SearchBlog(ctx *gin.Context) {
	searchParam := ctx.Query("search")
	var searchTitle string
	var searchContent string
	var searchTagId int
	if intValue, err := strconv.Atoi(searchParam); err == nil {
		log.Println(intValue, "from intvalue")
		searchTagId = intValue
	} else {
		searchContent = searchParam
		searchTitle = searchParam
	}
	log.Println(searchContent, " ", searchTitle, " ", searchTagId)
	blogSearchResponse := bh.blogService.SearchBlog(searchTitle, searchContent, searchTagId)
	if blogSearchResponse == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No blog found"})
		return
	}
	ctx.JSON(http.StatusFound, gin.H{"blogs": blogSearchResponse})
}

func ConvertId(userId interface{}) uint {
	var uid uint
	switch id := userId.(type) {
	case int:
		uid = uint(id)
	case float64:
		uid = uint(id)
	case uint:
		uid = id
	default:
		log.Fatalln("Unsupported type of id")
	}
	return uid
}
