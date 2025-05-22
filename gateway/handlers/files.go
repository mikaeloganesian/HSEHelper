package handlers

import (
	"gateway/config"
	"gateway/services"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFile godoc
// @Summary      Get file by ID
// @Description  Fetch a file from file-storing by its ID and stream it to the client
// @Tags         File Get
// @Accept       json
// @Produce      json
// @Param        id path string true "File ID"
// @Success      200 {file} string "File content"
// @Failure      400 {object} gin.H "Bad request"
// @Failure      500 {object} gin.H "Internal server error"
// @Router       /files/{id} [get]

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

// ListFiles godoc
// @Summary      List all files
// @Description  Fetch a list of all files from file-storing
// @Tags         File Get
// @Accept       json
// @Produce      json
// @Success      200 {array} services.FileResponse "List of files"
// @Failure      500 {object} gin.H "Internal server error"
// @Router       /files [get]

func ListFiles(c *gin.Context) {
	files, err := services.ListFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}
