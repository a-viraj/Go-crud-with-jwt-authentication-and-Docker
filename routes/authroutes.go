package routes

import (
	"github.com/a-viraj/project/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/user/signup", controller.Signup())
	r.POST("/user/login", controller.Login())
}
