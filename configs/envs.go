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
	GinMode              = os.Getenv("GIN_MODE")
	AppHost              = os.Getenv("BE_HOST")
	AppPort              = os.Getenv("BE_PORT")
	MariaUsername        = os.Getenv("MARIA_USERNAME")
	MariaPassword        = os.Getenv("MARIA_PASSWORD")
	MariaHost            = os.Getenv("MARIA_HOST")
	MariaPort            = os.Getenv("MARIA_PORT")
	AuthCoreDbName       = os.Getenv("AUTH_CORE_DB_NAME")
	MariaCharset         = os.Getenv("MARIA_CHARSET")
	MariaParseTime       = os.Getenv("MARIA_PARSE_TIME")
	MariaLoc             = os.Getenv("MARIA_LOC")
	MariaMaxOpenConns    = os.Getenv("MARIA_MAX_OPEN_CONNS")
	MariaMaxIdleConns    = os.Getenv("MARIA_MAX_IDLE_CONNS")
	MariaConnMaxLifetime = os.Getenv("MARIA_CONN_MAX_LIFETIME")
	AuthCoreDbDsn        = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		MariaUsername,
		MariaPassword,
		MariaHost,
		MariaPort,
		AuthCoreDbName,
		MariaCharset,
		MariaParseTime,
		MariaLoc,
	)
)

// SetupLogger setup logger
func SetupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// SetupDatabase setup database
func SetupDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", AuthCoreDbDsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	intMaxConn, err := strconv.Atoi(MariaMaxOpenConns)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MARIA_MAX_OPEN_CONNS to int: %w", err)
	}
	intMaxIdle, err := strconv.Atoi(MariaMaxIdleConns)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MARIA_MAX_IDLE_CONNS to int: %w", err)
	}
	intConnMaxLifetime, err := strconv.Atoi(MariaConnMaxLifetime)
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
