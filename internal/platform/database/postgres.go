package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	var dsn string

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		dsn = databaseURL
	} else if supabaseURL := os.Getenv("SUPABASE_DB_URL"); supabaseURL != "" {
		dsn = supabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=10",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
		)
		if strings.Contains(cfg.Host, "supabase.co") || strings.Contains(cfg.Host, "supabase.com") {
			dsn += " fallback_application_name=keep-your-house-clean"
		}
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewPostgresDBFromConfig() (*sql.DB, error) {
	host := getEnv("DB_HOST", getEnv("SUPABASE_DB_HOST", "localhost"))
	portStr := getEnv("DB_PORT", getEnv("SUPABASE_DB_PORT", "5432"))
	port, _ := strconv.Atoi(portStr)
	user := getEnv("DB_USER", getEnv("SUPABASE_DB_USER", "postgres"))
	password := getEnv("DB_PASSWORD", getEnv("SUPABASE_DB_PASSWORD", "postgres"))
	dbName := getEnv("DB_NAME", getEnv("SUPABASE_DB_NAME", "postgres"))
	sslMode := getEnv("DB_SSLMODE", getEnv("SUPABASE_DB_SSLMODE", ""))

	if sslMode == "" {
		if strings.Contains(host, "supabase.co") || strings.Contains(host, "supabase.com") {
			sslMode = "require"
		} else {
			sslMode = "disable"
		}
	}

	cfg := Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}

	return NewPostgresDB(cfg)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
