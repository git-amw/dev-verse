package routers

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/git-amw/backend/handlers"
	"gorm.io/gorm"

	"github.com/git-amw/backend/middleware"
	"github.com/git-amw/backend/services"
)

var DBInstance *gorm.DB
var ESClient *elasticsearch.Client

func SetupRouter() *gin.Engine {
	accountHandler := handlers.NewAccountHandler(services.NewAccountService(DBInstance))
	blogHandler := handlers.NewBlogHandler(services.NewBlogService(DBInstance, services.NewElasticSearchService(ESClient)))

	routes := gin.Default()
	routes.POST("/signup", accountHandler.CreateUser)
	routes.POST("/login", accountHandler.SignInUser)

	routes.GET("/allblogs", blogHandler.GetAllBlog)
	routes.POST("/createblog", middleware.AuthMiddleware, blogHandler.CreateBlog)
	routes.PATCH("/increaselikes/:id", blogHandler.IncreaseLike)
	routes.DELETE("/deleteblog/:id", blogHandler.DeleteBlog)
	routes.PUT("/updateblog", blogHandler.UpdateBlog)

	routes.GET("/alltags", blogHandler.GetAllTags)
	routes.GET("/searchtags", blogHandler.SearchTags)

	return routes
}
