package internal

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	_ "github.com/marcboeker/go-duckdb"
)

//go:embed db/migrations
var migrations embed.FS

type Note struct {
	UUID    string `json:"uuid"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Note    string `json:"notes"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

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

	u, err := uuid.NewV7()
	if err != nil {
		slog.Info("error creating UUID")
	}

	us, err := uuid.NewV7()
	if err != nil {
		slog.Info("error creating User UUID")
	}

	txn, err := db.Begin()
	if err != nil {
		slog.Info("error starting transactions")
	}

	sql := fmt.Sprintf("INSERT INTO NOTES (UUID, user, title, content, created_at, updated_at) VALUES ('%v', '%s' '%s', '%s', '%s', '%s')", u, us, title, note, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))
	slog.Info("SQL Query to debug", "query", sql)
	res, err := txn.ExecContext(ctx, sql)
	if err != nil {
		slog.Error("error commitng transaction, rolling back", err)
		txn.Rollback()
	}

	slog.Info("commit result", "message", res)
}

func GetAllNotes(ctx context.Context) Note {
	db, err := duckDB()
	note := Note{}
	if err != nil {
		slog.Info("error creating db")
	}

	res, err := db.QueryContext(ctx, "SELECT * FROM Notes")
	if err != nil {
		slog.Info("error starting transactions")
	}

	slog.Info("check for all the notes", "note", res)

	if res.Next() {
		slog.Info("check for all the notes", "note", res)
		if err := res.Scan(note.UUID, note.Title, note.Note, note.Created, note.Updated); err != nil {
			slog.Error("Error deserializing request")
		}
		return note
	}
	return Note{}
}
