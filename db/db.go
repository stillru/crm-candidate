package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "./crm.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Создаем таблицы, если они не существуют
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

// Exec выполняет SQL-запрос без возврата строк
func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

// QueryRow выполняет SQL-запрос и возвращает одну строку
func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

// Query выполняет SQL-запрос и возвращает несколько строк
func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS candidates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT,
			last_name TEXT,
			email TEXT,
			phone TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS vacancies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT,
			company_id INTEGER,
			FOREIGN KEY (company_id) REFERENCES companies(id)
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
