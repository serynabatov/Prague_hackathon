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
	keyManagement.GET("/get-private-key", controllers.GetPrivateKey)

	nfts := r.Group("/api/nfts")
	nfts.Use(middlewares.JwtAuthMiddleware())
	nfts.GET("", controllers.GetNfts)
	nfts.POST("/mint", controllers.MintNFT)

	collections := r.Group("/api/collections")
	collections.Use(middlewares.JwtAuthMiddleware())
	collections.POST("/setup", controllers.Setup)

	orgnizator_or_attendee := r.Group("/api/org_or_attendee")
	orgnizator_or_attendee.Use(middlewares.JwtAuthMiddleware())
	orgnizator_or_attendee.GET("", controllers.GetTheInfo)
	r.Run(":8080")
}
