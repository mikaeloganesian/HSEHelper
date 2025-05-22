package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway/config"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type FileFullResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Content   []byte `json:"content"`
	CreatedAt string `json:"created_at"`
}

// UploadFile загружает файл в file-storing
// и возвращает информацию о загруженном файле.
func UploadFile(fileHeader *multipart.FileHeader) (*FileResponse, error) {
	bodyReader, bodyWriter := io.Pipe()
	mw := multipart.NewWriter(bodyWriter)

	go func() {
		defer bodyWriter.Close()
		defer mw.Close()

		part, err := mw.CreateFormFile("file", fileHeader.Filename)
		if err != nil {
			return
		}

		src, err := fileHeader.Open()
		if err != nil {
			return
		}
		defer src.Close()

		io.Copy(part, src)
	}()

	req, err := http.NewRequest("POST", config.FileStoringURL+"/upload", bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("file-storing service returned non-200 status")
	}

	var result FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Получение файла по ID
func GetFileByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := http.Get(config.FileStoringURL + "/files/" + id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch file from storage"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "file-storing service returned non-200 status"})
		return
	}

	// Прокидываем заголовки от file-storing (например, Content-Type, Content-Disposition)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Устанавливаем статус
	c.Status(resp.StatusCode)

	// Прокидываем тело
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stream file content"})
	}
}

// Получение всех файлов
func ListFiles() ([]FileResponse, error) {
	resp, err := http.Get(config.FileStoringURL + "/files")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("file-storing service returned non-200 status")
	}

	var results []FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetFullFileByID(id string) ([]byte, error) {
	resp, err := http.Get(config.FileStoringURL + "/files/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("file-storing service returned non-200 status")
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func GetFileContentURL(id int) string {
	return fmt.Sprintf("%s/files/%d", config.FileStoringURL, id)
}
