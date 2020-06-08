// Package dbutil provides utilities for operating database.
package dbutil

import "database/sql"

// Connection represents a database connection.
type Connection interface {
	Open() (*sql.DB, error)
}

// Open opens a connection.
func Open(conn Connection, f func(db *sql.DB) error) (err error) {
	db, err := conn.Open()
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return
	}
	err = f(db)
	return
}
