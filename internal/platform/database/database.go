package database

import (
	"database/sql"
	"fmt"
)

func NewDatabase(path, driver string) (*sql.DB, error) {
	db, err := sql.Open(driver, path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %v", err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			address JSON,
			password TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS mail_validations (
			email TEXT PRIMARY KEY,
			code INTEGER NOT NULL,
			expired_at DATETIME NOT NULL
		);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %v", err)
	}

	return db, nil
}
