package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

const (
	ErrNotFound                   = DatabaseError("Could not get value from database, no such index")
	ErrCouldNotConnect            = DatabaseError("Could not connect to specified database file")
	ErrCouldNotCreateTable        = DatabaseError("Could not create table")
	ErrCouldNotQueryDatabase      = DatabaseError("Could not query database")
	ErrCouldNotExecuteQuery       = DatabaseError("Could not execute query to database")
	ErrCouldNotDeleteFromDatabase = DatabaseError("Could not execute delete query to database")
	ErrNoValueModified            = DatabaseError("No value was modified")
)

type DatabaseError string

func (e DatabaseError) Error() string {
	return string(e)
}

type Database struct {
	conn *sql.DB
}

func (db *Database) CloseConnection() {
	db.conn.Close()
}

func CreateDB(filePath string) *Database {
	db, err := sql.Open("sqlite", filePath)

	if err != nil {
		panic(err)
	}

	createIfNotExist(db)

	return &Database{db}
}

func createIfNotExist(db *sql.DB) {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS podcast (id INTEGER PRIMARY KEY, title TEXT, url TEXT)",
	)

	if err != nil {
		panic(err)
	}
}
