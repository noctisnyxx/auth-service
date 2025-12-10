package configs

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

var (
	GIN_MODE                = os.Getenv("GIN_MODE")
	BE_HOST                 = os.Getenv("BE_HOST")
	BE_PORT                 = os.Getenv("BE_PORT")
	MARIA_USERNAME          = os.Getenv("MARIA_USERNAME")
	MARIA_PASSWORD          = os.Getenv("MARIA_PASSWORD")
	MARIA_HOST              = os.Getenv("MARIA_HOST")
	MARIA_PORT              = os.Getenv("MARIA_PORT")
	AUTH_CORE_DB_NAME       = os.Getenv("AUTH_CORE_DB_NAME")
	MARIA_CHARSET           = os.Getenv("MARIA_CHARSET")
	MARIA_PARSE_TIME        = os.Getenv("MARIA_PARSE_TIME")
	MARIA_LOC               = os.Getenv("MARIA_LOC")
	MARIA_MAX_OPEN_CONNS    = os.Getenv("MARIA_MAX_OPEN_CONNS")
	MARIA_MAX_IDLE_CONNS    = os.Getenv("MARIA_MAX_IDLE_CONNS")
	MARIA_CONN_MAX_LIFETIME = os.Getenv("MARIA_CONN_MAX_LIFETIME")
	AUTH_CORE_DB_DSN        = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		MARIA_USERNAME,
		MARIA_PASSWORD,
		MARIA_HOST,
		MARIA_PORT,
		AUTH_CORE_DB_NAME,
		MARIA_CHARSET,
		MARIA_PARSE_TIME,
		MARIA_LOC,
	)
)

// setup logger
func SetupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// setup database
func SetupDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", AUTH_CORE_DB_DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	intMaxConn, err := strconv.Atoi(MARIA_MAX_OPEN_CONNS)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MARIA_MAX_OPEN_CONNS to int: %w", err)
	}
	intMaxIdle, err := strconv.Atoi(MARIA_MAX_IDLE_CONNS)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MARIA_MAX_IDLE_CONNS to int: %w", err)
	}
	intConnMaxLifetime, err := strconv.Atoi(MARIA_CONN_MAX_LIFETIME)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MARIA_CONN_MAX_LIFETIME to int: %w", err)

	}
	db.SetMaxOpenConns(intMaxConn)
	db.SetMaxIdleConns(intMaxIdle)
	db.SetConnMaxLifetime(time.Duration(intConnMaxLifetime) * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
