package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := os.Getenv("DATABASE_URL")
	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}

	if err := conn.Ping(ctx); err != nil {
		panic(fmt.Errorf("failed to ping DB: %w", err))
	}

	fmt.Println("Connected to DB successfully")
	DB = conn
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
