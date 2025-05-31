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
	public.POST("/generate-otp", controllers.GenerateOtp)
	public.GET("/google/login", controllers.GoogleLogin)
	public.GET("/google/callback", controllers.GoogleCallback)

	user := r.Group("/api/user")
	user.Use(middlewares.JwtAuthMiddleware())
	user.GET("", controllers.CurrentUser)

	keyManagement := r.Group("/api/key-management")
	keyManagement.Use(middlewares.JwtAuthMiddleware())
	keyManagement.GET("/get-private-key", controllers.GetPrivateKey)

	r.Run(":8080")
}
