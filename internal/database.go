package internal

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	_ "github.com/marcboeker/go-duckdb"
)

//go:embed db/migrations
var migrations embed.FS

func duckDB() (*sql.DB, error) {
	db, err := sql.Open("duckdb", "internal/db/embed.db")
	if err != nil {
		slog.Error("error connecting to the database", "message", err)
	}

	return db, nil
}

func CreateDB() {
	sql, err := migrations.ReadFile("db/migrations/001_init_schema.sql")
	if err != nil {
		slog.Info("error", err)
	}

	db, err := duckDB()
	if err != nil {
		slog.Info("error creating db")
	}

	db.Exec(string(sql))
}

func CreateNote(ctx context.Context, title, note string) {
	db, err := duckDB()
	if err != nil {
		slog.Info("error creating db")
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		slog.Info("error creating UUID")
	}

	txn, err := db.Begin()
	if err != nil {
		slog.Info("error starting transactions")
	}
	stmt, err := txn.PrepareContext(ctx, fmt.Sprintf("INSERT INTO NOTES (UUID, title, content) VALUES %v, %s, %s", uuid, title, note))
	if err != nil {
		slog.Info("error preparing context")
	}
	defer stmt.Close()

	if err := txn.Commit(); err != nil {
		slog.Error("error commitng transaction, rolling back")
		txn.Rollback()
	}
}
