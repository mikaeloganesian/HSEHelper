package handlers

import (
	"gateway/config"
	"gateway/services"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /files/:id
func GetFile(c *gin.Context) {
	id := c.Param("id")

	// Запрос к file-storing
	resp, err := http.Get(config.FileStoringURL + "/files/" + id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch file: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "File-storing service returned non-200 status"})
		return
	}

	// Прокидываем заголовки (особенно важно: Content-Type и Content-Disposition)
	for k, v := range resp.Header {
		for _, val := range v {
			c.Writer.Header().Add(k, val)
		}
	}

	c.Status(resp.StatusCode)

	// Стримим файл в ответ
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream file content"})
	}
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
