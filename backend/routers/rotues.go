package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/git-amw/backend/controllers"
	"github.com/git-amw/backend/middleware"
	"github.com/git-amw/backend/services"
)

var (
	accountController = controllers.NewAccountController(services.NewAccountService())
	blogController    = controllers.NewBlogController(services.NewBlogService())
)

func SetupRouter() *gin.Engine {

	routes := gin.Default()
	routes.POST("/signup", accountController.CreateUser)
	routes.POST("/login", accountController.SignInUser)

	routes.GET("/allblogs", blogController.GetAllBlog)
	routes.POST("/createblog", middleware.AuthMiddleware, blogController.CreateBlog)

	routes.GET("/alltags", blogController.GetAllTags)

	return routes
}
