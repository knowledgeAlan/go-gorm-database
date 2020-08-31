package routers

import (
	"vscode/go-gorm-database/controller"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {

	router := gin.Default()

	user := router.Group("user")
	{
		// user.POST("/register", controller.UserRegister)
		user.POST("/login", controller.Login01)
		user.POST("/todo", controller.CreateTodo)
		user.POST("/logout", controller.LoginOut)
		user.POST("/refresh", controller.RefreshToken)
	}
	return router

}
