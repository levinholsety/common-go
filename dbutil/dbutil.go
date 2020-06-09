// Package dbutil provides utilities for operating database.
package dbutil

import "database/sql"

// Connection represents a database connection.
type Connection interface {
	Open() (*sql.DB, error)
}

// Open creates a sql.DB instance and closes it after onOpened has been invoked.
func Open(conn Connection, onOpened func(db *sql.DB) error) (err error) {
	db, err := conn.Open()
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return
	}
	err = onOpened(db)
	return
}
