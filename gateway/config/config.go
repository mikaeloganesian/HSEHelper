package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	FileStoringURL  string
	FileAnalysisURL string
)

func LoadEnv() {
	err := godotenv.Load() // загружает переменные из .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	FileStoringURL = os.Getenv("FILE_STORING_URL")
	FileAnalysisURL = os.Getenv("FILE_ANALYSIS_URL")

	if FileStoringURL == "" || FileAnalysisURL == "" {
		log.Fatal("Missing FILE_STORING_URL or FILE_ANALYSIS_URL in environment variables")
	}
}
