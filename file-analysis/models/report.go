package models

import "time"

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
	ID         int
	FileName   string
	Paragraphs int
	Words      int
	Characters int
	Hash       string
	CreatedAt  time.Time
}
