package models

import (
	"context"
	"file-storing/config"
	"time"
)

type File struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Content   []byte    `json:"content"` // полный файл с контентом
	CreatedAt time.Time `json:"created_at"`
}

// Структура для списка файлов без контента
type FileSummary struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func InsertFile(file *File) error {
	query := `INSERT INTO files (name, content) VALUES ($1, $2) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.DB.QueryRow(ctx, query, file.Name, file.Content).Scan(&file.ID, &file.CreatedAt)
	return err
}

func GetFileByID(id int) (*File, error) {
	query := `SELECT id, name, content, created_at FROM files WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var f File
	err := config.DB.QueryRow(ctx, query, id).Scan(&f.ID, &f.Name, &f.Content, &f.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func ListAllFiles() ([]FileSummary, error) {
	query := `SELECT id, name, created_at FROM files ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []FileSummary
	for rows.Next() {
		var f FileSummary
		err := rows.Scan(&f.ID, &f.Name, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}
