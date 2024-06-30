package main

import (
	"github.com/a-viraj/project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := "8000"
	router := gin.New()
	router.Use(gin.Logger())
	router.GET("api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "access granted"})
	})
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	router.Run(":" + port)
}
