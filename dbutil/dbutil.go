// Package dbutil provides utilities for operating database.
package dbutil

import (
	"bytes"
	"database/sql"
	"strings"
)

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

// Split split sql text with semicolon.
func Split(sqlText string) (result []string) {
	var (
		inSingleQuotes = false
		inDoubleQuotes = false
		escape         = false
	)
	buf := &bytes.Buffer{}
	appendStr := func() {
		text := buf.String()
		buf.Reset()
		text = strings.TrimSpace(text)
		if len(text) > 0 {
			result = append(result, text)
		}
	}
	for _, r := range sqlText {
		if inSingleQuotes {
			if escape {
				escape = false
			} else if r == '\\' {
				escape = true
			} else if r == '\'' {
				inSingleQuotes = false
			}
		} else if inDoubleQuotes {
			if escape {
				escape = false
			} else if r == '\\' {
				escape = true
			} else if r == '"' {
				inDoubleQuotes = false
			}
		} else {
			if r == '\'' {
				inSingleQuotes = true
			} else if r == '"' {
				inDoubleQuotes = true
			} else if r == ';' {
				appendStr()
				continue
			}
		}
		buf.WriteRune(r)
	}
	appendStr()
	return
}
