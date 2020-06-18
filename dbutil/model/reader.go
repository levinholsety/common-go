package model

import (
	"database/sql"
	"errors"
)

// Reader provides methods to read model from database.
type Reader interface {
	ReadSchemas(db *sql.DB, m *Model) error
	ReadTables(db *sql.DB, schema *Schema) error
	ReadTable(db *sql.DB, schemaName string, table *Table) error
	ReadColumns(db *sql.DB, schemaName string, table *Table) error
}

// Errors
var (
	ErrTableNotFound = errors.New("table not found")
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
	err = r.ReadTable(db, schemaName, result)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrTableNotFound
		}
		return
	}
	err = r.ReadColumns(db, schemaName, result)
	if err != nil {
		return
	}
	return
}
