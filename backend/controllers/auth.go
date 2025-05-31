package controllers

import (
	"backend/cloud"
	"backend/models"
	"backend/utils/password"
	"backend/utils/token"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Email = input.Email
	u.Password = input.Password
	u.Authenicated = false

	_, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := cloud.SetUpTotp(u.ID, u.Email)

	if err != nil {
		log.Fatal("[Error] generating the totp")
		return
	}

	fmt.Println("Scan this QR/URL in Google Authenticator or Authy:", key)
	fmt.Println("Then use a URL like: http://localhost:8080/get-private-key?user=user123&code=XXXXXX")

	c.JSON(http.StatusOK, gin.H{"otp": key})
}

func GenerateOtp(c *gin.Context) {
	email := c.Query("email")

	u, err := models.GetByEmail(email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := cloud.SetUpTotp(u.ID, u.Email)

	if err != nil {
		log.Fatal("[Error] generating the totp")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate a secure password"})
		return
	}

	fmt.Println("Scan this QR/URL in Google Authenticator or Authy:", key)
	fmt.Println("Then use a URL like: http://localhost:8080/get-private-key?user=user123&code=XXXXXX")

	c.JSON(http.StatusOK, gin.H{"otp": key})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Email = input.Email
	u.Password = input.Password

	token, err := models.LoginCheck(u.Email, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL: "http://localhost:8080/api/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

const oauthStateString = "random-secret"

func GoogleLogin(c *gin.Context) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	googleOauthConfig.ClientID = clientId
	googleOauthConfig.ClientSecret = clientSecret
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	googleOauthConfig.ClientID = clientId
	googleOauthConfig.ClientSecret = clientSecret
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	tokenGoogle, err := googleOauthConfig.Exchange(context.Background(), code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), tokenGoogle)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
		return
	}

	user := models.User{}

	user, err = models.GetByEmail(userInfo.Email)

	if err != nil {
		password, err := password.GenerateRandomPassword(24)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate a secure password"})
			return
		}
		user.Email = userInfo.Email
		user.Password = string(password)
		user.Authenicated = false
		user.SaveUser()
		key, err := cloud.SetUpTotp(user.ID, user.Email)

		if err != nil {
			log.Fatal("[Error] generating the totp")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
			return
		}

		fmt.Println("Scan this QR/URL in Google Authenticator or Authy:", key)
		fmt.Println("Then use a URL like: http://localhost:8080/get-private-key?user=user123&code=XXXXXX")

		c.JSON(http.StatusOK, gin.H{"otp": key})
		return
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate a token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
