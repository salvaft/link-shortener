package persistance

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
	if err != nil {
		log.Fatalf("%-20s fatal error initializing db: %v\n", "NewSQLite", err)
	}
	return newDB
}

func (db *SQLiteDB) Close() {
	db.db.Close()
}
