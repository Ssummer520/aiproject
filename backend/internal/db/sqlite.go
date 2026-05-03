package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const DefaultPath = "data/travel.db"

var (
	shared  *sql.DB
	once    sync.Once
	initErr error
)

func Open() (*sql.DB, error) {
	once.Do(func() {
		path := os.Getenv("TRAVEL_DB_PATH")
		if path == "" {
			path = DefaultPath
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			initErr = err
			return
		}
		db, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_busy_timeout=5000")
		if err != nil {
			initErr = err
			return
		}
		db.SetMaxOpenConns(1)
		if err := migrate(db); err != nil {
			_ = db.Close()
			initErr = err
			return
		}
		shared = db
	})
	return shared, initErr
}

func migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS tokens (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			expires_at TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_user_type ON tokens(user_id, type);`,
		`CREATE TABLE IF NOT EXISTS favorites (
			user_id TEXT NOT NULL,
			destination_id INTEGER NOT NULL,
			is_favorite INTEGER NOT NULL DEFAULT 1,
			updated_at TEXT NOT NULL,
			PRIMARY KEY(user_id, destination_id)
		);`,
		`CREATE TABLE IF NOT EXISTS history (
			user_id TEXT NOT NULL,
			destination_id INTEGER NOT NULL,
			viewed_at TEXT NOT NULL,
			PRIMARY KEY(user_id, destination_id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_history_user_viewed ON history(user_id, viewed_at DESC);`,
		`CREATE TABLE IF NOT EXISTS stats (
			destination_id INTEGER PRIMARY KEY,
			views INTEGER NOT NULL DEFAULT 0,
			favorites INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			destination_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			city TEXT NOT NULL,
			cover TEXT NOT NULL,
			check_in TEXT NOT NULL,
			check_out TEXT NOT NULL,
			guests INTEGER NOT NULL,
			total_price REAL NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			cancelled_at TEXT,
			PRIMARY KEY(user_id, id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_user_created ON bookings(user_id, created_at DESC);`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			message TEXT NOT NULL,
			link TEXT NOT NULL,
			read INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			PRIMARY KEY(user_id, id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}
