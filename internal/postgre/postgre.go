package postgre

import (
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresDBConfig struct {
	Username    string
	Password    string
	Host        string
	Database    string
	MaxConn     int
	MaxConnIdle int
}

func NewPostgresDB(config PostgresDBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		url.QueryEscape(config.Username),
		url.QueryEscape(config.Password),
		config.Host,
		config.Database,
	)

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres connection config failed: %w", err)
	}

	db := sqlx.NewDb(stdlib.OpenDB(*connConfig), "pgx")
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping to database failed: %w", err)
	}

	db.SetMaxOpenConns(config.MaxConn)
	db.SetMaxIdleConns(config.MaxConnIdle)

	return db, nil
}
