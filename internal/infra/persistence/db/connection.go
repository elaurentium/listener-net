package db

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
)


type DBConnection struct {
	User	 string
	Password string
	DBName   string
	Host     string
	Port     string
}


func NewDBConnection() (*sql.DB, error) {
	config := &DBConnection{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}