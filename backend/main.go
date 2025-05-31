package main

import (
	"backend/cloud"
	"backend/controllers"
	"backend/middlewares"
	"backend/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDataBase()
	key, err := cloud.SetUpTotp(1, "sergio.nabatini@prague.cz")

	if err != nil {
		log.Fatal("[Error] generating the totp")
		return
	}

	fmt.Println("Scan this QR/URL in Google Authenticator or Authy:", key)
	fmt.Println("Then use a URL like: http://localhost:8080/get-private-key?user=user123&code=XXXXXX")

	r := gin.Default()

	public := r.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
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
