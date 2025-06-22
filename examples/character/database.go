package main

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS characters (
		id INTEGER PRIMARY KEY,
		character_data TEXT NOT NULL,
		cached_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_characters_id ON characters(id);
	CREATE INDEX IF NOT EXISTS idx_characters_cached_at ON characters(cached_at);
	`

	_, err := d.db.Exec(query)
	return err
}

func (d *Database) GetCharacter(id uint32) ([]byte, bool, error) {
	var data string
	var cachedAt time.Time

	query := "SELECT character_data, cached_at FROM characters WHERE id = ?"
	err := d.db.QueryRow(query, id).Scan(&data, &cachedAt)
	
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	// 檢查快取是否過期 (24小時)
	if time.Since(cachedAt) > 24*time.Hour {
		return nil, false, nil
	}

	return []byte(data), true, nil
}

func (d *Database) SaveCharacter(id uint32, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	now := time.Now()
	query := `
	INSERT OR REPLACE INTO characters (id, character_data, cached_at, updated_at) 
	VALUES (?, ?, ?, ?)
	`

	_, err = d.db.Exec(query, id, string(jsonData), now, now)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}