package handlers

import (
	"gateway/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /files/:id
func GetFile(c *gin.Context) {
	id := c.Param("id")

	result, err := services.GetFileByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch file" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GET /files
func ListFiles(c *gin.Context) {
	files, err := services.ListFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}
