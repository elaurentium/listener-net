package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)


type DBConnection struct {
	User	 string
	Password string
	DBName   string
}


func NewDBConnection() (*pgxpool.Pool, error) {
	config := &DBConnection{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
	}

	dns := fmt.Sprintf("%s:%s@/%s", config.User, config.Password, config.DBName)

	db, err := pgxpool.Connect(context.Background(), dns)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}