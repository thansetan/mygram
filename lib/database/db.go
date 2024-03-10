package database

import (
	"database/sql"
	_ "embed"
	"final-project/lib/config"
	"fmt"
)

func New(conf config.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("database.New: %w", err)
	}

	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)
	db.SetConnMaxIdleTime(conf.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database.New: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("database.New: %w", err)
	}

	return db, nil
}

//go:embed db.sql
var q string

func createTables(db *sql.DB) error {
	_, err := db.Exec(q)
	return err
}
