package handlers

import (
	"gateway/services"
	"io"
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

func UploadAndAnalyze(c *gin.Context) {
	// 1. Получение файла от пользователя
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// 2. Загрузка в file-storing
	uploadResp, err := services.UploadFile(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// 3. Получение байтового потока (контента) по ID
	contentResp, err := http.Get(services.GetFileContentURL(uploadResp.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file content: " + err.Error()})
		return
	}
	defer contentResp.Body.Close()

	// 4. Чтение потока байт -> в строку
	data, err := io.ReadAll(contentResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content: " + err.Error()})
		return
	}

	// 5. Анализ текста
	analysisReq := services.AnalyzeRequest{
		Text:     string(data),
		FileName: uploadResp.Name,
	}
	analysisRes, err := services.AnalyzeText(analysisReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze text: " + err.Error()})
		return
	}

	// 6. Ответ клиенту
	c.JSON(http.StatusOK, gin.H{
		"file_id":    uploadResp.ID,
		"file_name":  uploadResp.Name,
		"created_at": uploadResp.CreatedAt,
		"analysis":   analysisRes,
	})
}
