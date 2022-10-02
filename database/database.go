package database

import (
	"database/sql"
)

type Database struct {
	dialect       string
	connectionStr string
}

func NewDatabase(dialect string, connectionStr string) *Database {
	return &Database{
		dialect:       dialect,
		connectionStr: connectionStr,
	}
}

/*
Open will open a connection to the
database.

NB: Remember to close the connection
when you are done!
*/
func (db *Database) Open() (*sql.DB, error) {
	return sql.Open(db.dialect, db.connectionStr)
}

func (db *Database) ConnectionString() string {
	return db.connectionStr
}
