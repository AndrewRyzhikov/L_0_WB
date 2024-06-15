package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"L_0_WB/internal/config"
)

type Storage struct {
	*sql.DB
}

func NewStorage(cfg config.PostgresConfig) (*Storage, error) {
	db, err := sql.Open("postgres", dataToPSQLConnection(cfg.Port, cfg.Host, cfg.User, cfg.Password, cfg.DbName))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}
	return &Storage{db}, nil
}

func dataToPSQLConnection(port int, host, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func (d *Storage) ExecuteInsert(ctx context.Context, query string, args ...interface{}) error {
	tx, err := d.Begin()
	if err != nil {
		return fmt.Errorf("transaction cannot start: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("transaction cannot execute and rollback: %w, %w", err, errRollback)
		}
		return fmt.Errorf("transaction cannot execute: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("transaction cannot commit: %w", err)
	}

	return nil
}
