package model

// StatementGenerator provides methods to generate statements.
type StatementGenerator interface {
	GenerateCreateTableStatement(table *Table) string
	GenerateAlterTableCommentStatement(table *Table) string
	GenerateDropTableStatement(tableName string) string
	GenerateAddColumnStatement(table *Table, columnIndex int) string
	GenerateModifyColumnStatement(table *Table, columnIndex int) string
	GenerateDropColumnStatement(tableName, columnName string) string
	GenerateAddPrimaryKeyStatement(table *Table) string
	GenerateDropPrimaryKeyStatement(tableName string) string
}
