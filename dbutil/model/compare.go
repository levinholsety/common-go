package model

import (
	"github.com/levinholsety/common-go/comm"
)

// CompareSchema compares the database schema with its old version and generates statement for altering database.
func CompareSchema(schema, oldSchema *Schema, onResult func(result ComparisonResult)) {
	if schema == nil || oldSchema == nil {
		return
	}
	for _, table := range schema.Tables {
		oldTable := findTable(oldSchema, table.Name)
		CompareTable(table, oldTable, onResult)
	}
	for _, oldTable := range oldSchema.Tables {
		table := findTable(schema, oldTable.Name)
		if table == nil {
			onResult(&TableRedundantComparisonResult{TableName: oldTable.Name})
		}
	}
	return
}

// CompareTable compares the table with its old version and generates statement for altering database.
func CompareTable(table, oldTable *Table, onResult func(result ComparisonResult)) {
	if table == nil && oldTable == nil {
		return
	}
	if table != nil && oldTable == nil {
		onResult(&TableMissingComparisonResult{Table: table})
		return
	}
	if table == nil && oldTable != nil {
		onResult(&TableRedundantComparisonResult{TableName: oldTable.Name})
		return
	}
	if table.Name != oldTable.Name {
		onResult(&TableChangedComparisonResult{Table: table, OldTable: oldTable})
		return
	}
	if table.Comment != oldTable.Comment {
		onResult(&TableCommentChangedComparisonResult{Table: table})
	}
	for i, column := range table.Columns {
		oldColumn := findColumn(oldTable, column.Name)
		if oldColumn == nil {
			onResult(&ColumnMissingComparisonResult{Table: table, ColumnIndex: i})
		} else {
			if !columnEqual(column, oldColumn) {
				onResult(&ColumnChangedComparisonResult{Table: table, ColumnIndex: i})
			}
		}
	}
	for _, oldColumn := range oldTable.Columns {
		column := findColumn(table, oldColumn.Name)
		if column == nil {
			onResult(&ColumnRedundantComparisonResult{TableName: table.Name, ColumnName: oldColumn.Name})
		}
	}
	pkColumns := table.PrimaryKeyColumnNames()
	oldPKColumns := oldTable.PrimaryKeyColumnNames()
	if len(pkColumns) > 0 && len(oldPKColumns) == 0 {
		onResult(&PrimaryKeyMissingComparisonResult{Table: table})
	} else if len(pkColumns) == 0 && len(oldPKColumns) > 0 {
		onResult(&PrimaryKeyRedundantComparisonResult{TableName: table.Name})
	} else if len(pkColumns) > 0 && len(oldPKColumns) > 0 && !comm.StringArrayEqual(pkColumns, oldPKColumns) {
		onResult(&PrimaryKeyChangedComparisonResult{Table: table})
	}
	return
}

func findTable(schema *Schema, tableName string) *Table {
	for _, tbl := range schema.Tables {
		if tbl.Name == tableName {
			return tbl
		}
	}
	return nil
}

func findColumn(table *Table, columnName string) *Column {
	for _, col := range table.Columns {
		if col.Name == columnName {
			return col
		}
	}
	return nil
}

func columnEqual(column1, column2 *Column) bool {
	return column1.Name == column2.Name &&
		column1.Type == column2.Type &&
		column1.Nullable == column2.Nullable &&
		column1.Comment == column2.Comment
}
