package routes

import (
	"github.com/a-viraj/project/controller"
	"github.com/a-viraj/project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.Use(middleware.Authenticate())
	r.GET("/users", controller.GetUsers())
	r.GET("/users/:userId", controller.GetUser())
}
