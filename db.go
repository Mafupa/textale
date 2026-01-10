package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func initDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "./data/textale.db")
	if err != nil {
		return nil, err
	}

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS channels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		retention_seconds INTEGER DEFAULT 300,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS banned_users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT NOT NULL,
		username TEXT,
		reason TEXT,
		banned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_preferences (
		ip_address TEXT PRIMARY KEY,
		username TEXT,
		last_channel TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	INSERT OR IGNORE INTO channels (name) VALUES ('general'), ('random'), ('admin');
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) IsBanned(ipAddress string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM banned_users WHERE ip_address = ?", ipAddress).Scan(&count)
	return count > 0, err
}

func (db *DB) GetChannels() ([]string, error) {
	rows, err := db.Query("SELECT name FROM channels ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		channels = append(channels, name)
	}
	return channels, nil
}

func (db *DB) SavePreference(ipAddress, username, channel string) error {
	_, err := db.Exec(`
		INSERT INTO user_preferences (ip_address, username, last_channel, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(ip_address) DO UPDATE SET
			username = excluded.username,
			last_channel = excluded.last_channel,
			updated_at = CURRENT_TIMESTAMP
	`, ipAddress, username, channel)
	return err
}
