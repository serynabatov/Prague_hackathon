package controllers

import (
	"backend/custom_flow"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTheInfo(c *gin.Context) {
	records, err := custom_flow.GetOrganizatorOrAttendee()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}
