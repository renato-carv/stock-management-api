package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Username string
	Password string
	Database string
	Port     string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Port:     os.Getenv("DB_PORT"),
	}
}

func SetupDb(cfg *DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port,
	)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = dbConn.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		
		CREATE TABLE IF NOT EXISTS users (
			ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			NAME TEXT NOT NULL,
			USERNAME VARCHAR(100) UNIQUE NOT NULL, 
			EMAIL VARCHAR(100) UNIQUE NOT NULL,
			PASSWORD TEXT NOT NULL,
			CREATED_AT TIMESTAMP DEFAULT now(),
			UPDATED_AT TIMESTAMP DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS stock (
			ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			NAME TEXT NOT NULL,
			QUANTITY INTEGER NOT NULL,
			CREATED_AT TIMESTAMP DEFAULT now(),
			UPDATED_AT TIMESTAMP DEFAULT now(),
			CREATED_BY UUID REFERENCES users(ID)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	log.Println("Table created successfully")
	return dbConn, nil
}
