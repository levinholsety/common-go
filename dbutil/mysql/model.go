package mysql

import (
	"database/sql"
	"strings"

	"github.com/levinholsety/common-go/dbutil/model"
)

// ReadModel reads database model of specified schema from database.
func ReadModel(db *sql.DB, schemaName string) (result *model.Model, err error) {
	schema := &model.Schema{Name: schemaName}
	if err = readTables(db, schema); err != nil {
		return
	}
	for _, table := range schema.Tables {
		if err = readColumns(db, schema, table); err != nil {
			return
		}
	}
	result = &model.Model{Schemas: []*model.Schema{schema}}
	return
}

func readSchemas(db *sql.DB, m *model.Model) (err error) {
	rows, err := db.Query(`select schema_name from information_schema.schemata`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		schema := &model.Schema{}
		if err = rows.Scan(&schema.Name); err != nil {
			return
		}
		m.Schemas = append(m.Schemas, schema)
	}
	return
}

func readTables(db *sql.DB, schema *model.Schema) (err error) {
	rows, err := db.Query(`select table_name,table_comment from information_schema.tables where table_schema = ? and table_type = 'BASE TABLE'`, schema.Name)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		table := &model.Table{}
		if rows.Scan(&table.Name, &table.Comment); err != nil {
			return
		}
		schema.Tables = append(schema.Tables, table)
	}
	return
}

func readColumns(db *sql.DB, schema *model.Schema, table *model.Table) (err error) {
	rows, err := db.Query(`select column_name,data_type,column_type,is_nullable,column_comment,column_key
from information_schema.columns
where table_schema = ? and table_name = ? order by ordinal_position`, schema.Name, table.Name)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		column := &model.Column{}
		var (
			nullable  string
			columnKey string
		)
		if err = rows.Scan(&column.Name, &column.DataType, &column.Type, &nullable, &column.Comment, &columnKey); err != nil {
			return
		}
		column.DataType = strings.ToUpper(column.DataType)
		switch column.DataType {
		case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT", "JSON":
			column.DataClass = model.Text
		case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT", "BIT", "FLOAT", "DOUBLE", "DECIMAL":
			column.DataClass = model.Number
		case "BINARY", "VARBINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
			column.DataClass = model.Binary
		case "DATE", "TIME", "YEAR", "DATETIME", "TIMESTAMP":
			column.DataClass = model.Time
		}
		column.Type = strings.ToUpper(column.Type)
		column.Nullable = nullable != "NO"
		column.IsPrimaryKey = columnKey == "PRI"
		table.Columns = append(table.Columns, column)
	}
	return
}
