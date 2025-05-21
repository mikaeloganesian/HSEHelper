package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"gateway/config"
	"net/http"
)

type AnalyzeRequest struct {
	Text     string `json:"text"`
	FileName string `json:"file_name"`
}

type AnalyzeResponse struct {
	Paragraphs    int  `json:"paragraphs"`
	Words         int  `json:"words"`
	Characters    int  `json:"characters"`
	IsPlagiarized bool `json:"is_plagiarized"`
}

type Report struct {
	ID         int    `json:"id"`
	FileName   string `json:"file_name"`
	Paragraphs int    `json:"paragraphs"`
	Words      int    `json:"words"`
	Characters int    `json:"characters"`
	Hash       string `json:"hash"`
	CreatedAt  string `json:"created_at"`
}

func AnalyzeText(req AnalyzeRequest) (*AnalyzeResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(config.FileAnalysisURL+"/analyze", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("analysis service returned non-200 status")
	}

	var result AnalyzeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetReportByID(id string) (*Report, error) {
	resp, err := http.Get(config.FileAnalysisURL + "/reports/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("file-analysis service returned non-200 status")
	}

	var report Report
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, err
	}

	return &report, nil
}
