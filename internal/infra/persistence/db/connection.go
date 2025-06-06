package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)


type DBConfig struct {
	Host		string
	Port		string
	User		string
	Password	string
	DBName		string
	SSLMode		string
}


func NewDBConnection() (*pgxpool.Pool, error) {
	config := &DBConfig{
		User: 		os.Getenv("DB_USER"),
		Password: 	os.Getenv("DB_PASSWORD"),
		DBName: 	os.Getenv("DB_NAME"),
		Host: 		os.Getenv("DB_HOST"),
		Port: 		os.Getenv("DB_PORT"),
		SSLMode: 	"disable",
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	
	db, err := pgxpool.Connect(context.Background(), dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}