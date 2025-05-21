package handlers

import (
	"gateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required" + err.Error()})
		return
	}

	result, err := services.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
