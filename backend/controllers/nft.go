package controllers

import (
	"backend/custom_flow"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNfts(c *gin.Context) {
	nfts, err := custom_flow.GetNfts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store private key"})
		return
	}

	c.JSON(http.StatusOK, nfts)
}

func MintNFT(c *gin.Context) {
	recipientAddress := c.Query("address") // 0x51846a0f69492bba

	var input custom_flow.NFTData

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	custom_flow.Mint(recipientAddress, input)

	c.JSON(http.StatusOK, gin.H{"message": "Ok"})
}
