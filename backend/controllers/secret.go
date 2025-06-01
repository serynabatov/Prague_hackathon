package controllers

import (
	"backend/cloud"
	"backend/custom_flow"
	"backend/models"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

func GetPrivateKey(c *gin.Context) {
	ctx := c.Request.Context()
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	email := c.Query("email")

	u, err := models.GetByEmail(email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := u.ID

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1) Validate TOTP code
	code := c.Query("code")
	if !cloud.ValidateTOTP(userID, code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid 2FA code"})
		return
	}

	client, err := cloud.NewSecretManagerClient(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot initialize secretmanager client"})
		return
	}
	defer client.Close()

	_, err = cloud.GetUserPrivateKey(ctx, client, projectID, userID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"user":      userID,
			"address":   "0x00023",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}
	log.Println(err)
	st := status.Convert(err)

	if st.Code() == codes.NotFound {
		if err := cloud.CreateUserSecret(ctx, client, projectID, userID); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Secret Manager secret"})
			return
		}

		flowAccount, genErr := custom_flow.CreateAccount()

		if genErr != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate private key"})
			return
		}

		if err := cloud.AddUserSecretVersion(ctx, client, projectID, userID, []byte(flowAccount.PrivateKey)); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store private key"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":      userID,
			"address":   flowAccount.Address,
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return

	}

	log.Printf("unexpected error accessing secret: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to retrieve private key"})
}

func generateRandomKey(nBytes int) (string, error) {
	data := make([]byte, nBytes)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}
