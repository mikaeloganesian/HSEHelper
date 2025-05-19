package handlers

import (
	"net/http"
	"strconv"

	"file-storing/models"

	"github.com/gin-gonic/gin"
)

func GetFileByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	file, err := models.GetFileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+file.Name)
	c.Data(http.StatusOK, "application/octet-stream", file.Content)
}

func ListFiles(c *gin.Context) {
	files, err := models.ListAllFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	c.JSON(http.StatusOK, files)
}
