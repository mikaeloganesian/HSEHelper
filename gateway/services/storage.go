package services

import (
	"encoding/json"
	"errors"
	"gateway/config"
	"io"
	"mime/multipart"
	"net/http"
)

type FileResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

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
func GetFileByID(id string) (*FileResponse, error) {
	resp, err := http.Get(config.FileStoringURL + "/files/" + id)
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
