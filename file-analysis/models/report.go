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
	ID         int       `json:"id" db:"id"`
	FileName   string    `json:"file_name" db:"file_name"`
	Paragraphs int       `json:"paragraphs" db:"paragraphs"`
	Words      int       `json:"words" db:"words"`
	Characters int       `json:"characters" db:"characters"`
	Hash       string    `json:"hash" db:"hash"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
