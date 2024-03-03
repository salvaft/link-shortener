package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLite(dbName string) *SQLiteDB {
	db, err := sql.Open("sqlite3", "file:"+dbName)
	if err != nil {
		log.Fatalf("%-20s fatal error connecting to db: %v\n", "NewSQLite", err)
	}
	newDB := &SQLiteDB{db}
	err = newDB.Init()
	if err != nil {
		log.Fatalf("%-20s fatal error initializing db: %v\n", "NewSQLite", err)
	}
	return newDB
}

func (db *SQLiteDB) Close() {
	db.db.Close()
}

func (db *SQLiteDB) Init() error {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS links (id INTEGER PRIMARY KEY, url TEXT)")
	if err != nil {
		return err
	}

	_, err = db.db.Exec("CREATE INDEX IF NOT EXISTS idx_url ON links(url)")

	if err != nil {
		return err
	}
	return nil
}
