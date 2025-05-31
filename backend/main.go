package main

import (
	"backend/controllers"
	"backend/middlewares"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDataBase()

	r := gin.Default()

	public := r.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.GET("/google/login", controllers.GoogleLogin)
	public.GET("/google/callback", controllers.GoogleCallback)

	user := r.Group("/api/user")
	user.Use(middlewares.JwtAuthMiddleware())
	user.GET("", controllers.CurrentUser)

	r.Run(":8080")
}
