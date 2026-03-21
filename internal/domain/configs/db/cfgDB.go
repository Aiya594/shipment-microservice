package configs

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type ConfigDB struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SSLMode  string
}

func NewConfigsDB() *ConfigDB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbMode := os.Getenv("DB_SSLMODE")

	return &ConfigDB{
		Host:     dbHost,
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Port:     dbPort,
		SSLMode:  dbMode,
	}
}

func (d *ConfigDB) Connect(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", d.ConnectionStr())
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}

func (d *ConfigDB) ConnectionStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}
