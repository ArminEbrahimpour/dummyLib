package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB(ctx context.Context, dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("init database failed: %w", err)
	}

	// enable wal mode for concurrent reading
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("enebaling wal mode failed: %w", err)
	}

	if err := ensureSchema(ctx, db); err != nil {
		return nil, fmt.Errorf("ensure schema failed: %w", err)
	}
	return db, nil
}

func ensureSchema(ctx context.Context, db *sql.DB) error {

	schema := `
	CREATE	TABLE IF NOT EXISTS books (
		id	TEXT PRIMARY KEY, 
		path TEXT NOT NULL UNIQUE,
		title TEXT,
		page_count INTEGER,
		last_page INTEGER DEFAULT 1,
		last_opened DATETIME,
		cover_path TEXT	
	);
	`
	_, err := db.ExecContext(ctx, schema)
	return err

}
