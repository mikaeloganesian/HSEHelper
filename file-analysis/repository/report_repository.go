package repository

import (
	"context"
	"fmt"
	"time"

	"file-analysis/config"
	"file-analysis/models"
)

func FindReportByID(id int) (*models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, file_name, paragraphs, words, characters, hash, created_at FROM reports WHERE id=$1`
	row := config.DB.QueryRow(ctx, query, id)

	var report models.Report
	err := row.Scan(&report.ID, &report.FileName, &report.Paragraphs, &report.Words, &report.Characters, &report.Hash, &report.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func FindReportByHash(hash string) (*models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, file_name, paragraphs, words, characters, hash, created_at FROM reports WHERE hash=$1`
	row := config.DB.QueryRow(ctx, query, hash)

	var report models.Report
	err := row.Scan(&report.ID, &report.FileName, &report.Paragraphs, &report.Words, &report.Characters, &report.Hash, &report.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func InsertReport(report *models.Report) error {
	existing, err := FindReportByHash(report.Hash)
	if err == nil && existing != nil {
		return fmt.Errorf("report with identical hash already exists, id=%d", existing.ID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO reports (file_name, paragraphs, words, characters, hash) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;`

	row := config.DB.QueryRow(ctx, query, report.FileName, report.Paragraphs, report.Words, report.Characters, report.Hash)
	return row.Scan(&report.ID, &report.CreatedAt)
}
