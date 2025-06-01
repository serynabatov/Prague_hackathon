package controllers

import (
	"backend/cloud"
	"backend/custom_flow"
	"backend/utils/token"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(c *gin.Context) {
	ctx := c.Request.Context()
	projectID := os.Getenv("GOOGLE_PROJECT_ID")

	userID, err := token.ExtractTokenID(c)

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

	keyBytes, err := cloud.GetUserPrivateKey(ctx, client, projectID, userID)
	if err == nil {

		custom_flow.Setup(string(keyBytes))
		c.JSON(http.StatusOK, gin.H{
			"user":      userID,
			"private":   string(keyBytes),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}
	log.Printf("unexpected error accessing secret: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to retrieve private key"})
}
