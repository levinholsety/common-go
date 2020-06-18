package model

import (
	"database/sql"
	"errors"
)

// Reader provides methods to read model from database.
type Reader interface {
	ReadSchemas(db *sql.DB, m *Model) error
	ReadTables(db *sql.DB, schema *Schema) error
	ReadColumns(db *sql.DB, schemaName string, table *Table) error
}

// Errors
var (
	ErrNoTables  = errors.New("no tables")
	ErrNoColumns = errors.New("no columns")
)

// ReadModel reads model of database.
func ReadModel(db *sql.DB, r Reader) (result *Model, err error) {
	result = &Model{}
	if err = r.ReadSchemas(db, result); err != nil {
		return
	}
	for _, schema := range result.Schemas {
		if err = r.ReadTables(db, schema); err != nil {
			return
		}
		for _, table := range schema.Tables {
			if err = r.ReadColumns(db, schema.Name, table); err != nil {
				return
			}
		}
	}
	return
}

// ReadSchema reads model of database schema from database.
func ReadSchema(db *sql.DB, schemaName string, r Reader) (result *Schema, err error) {
	result = &Schema{Name: schemaName}
	if err = r.ReadTables(db, result); err != nil {
		return
	}
	if len(result.Tables) == 0 {
		err = ErrNoTables
		return
	}
	for _, table := range result.Tables {
		if err = r.ReadColumns(db, schemaName, table); err != nil {
			return
		}
	}
	return
}

// ReadTable reads model of database table from database.
func ReadTable(db *sql.DB, schemaName, tableName string, r Reader) (result *Table, err error) {
	result = &Table{Name: tableName}
	err = r.ReadColumns(db, schemaName, result)
	if err != nil {
		return
	}
	if len(result.Columns) == 0 {
		err = ErrNoColumns
		return
	}
	return
}
