package services

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"

	"file-analysis/models"
)

func CalculateHash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func AnalyzeText(text string) models.AnalyzeResponse {
	paragraphs := strings.Count(text, "\n\n") + 1
	words := len(strings.Fields(text))
	characters := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			characters++
		}
	}

	return models.AnalyzeResponse{
		Paragraphs: paragraphs,
		Words:      words,
		Characters: characters,
	}
}
