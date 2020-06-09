package model

import (
	"io"

	"github.com/levinholsety/common-go/comm"
)

// ComparisonResults represents comparison results.
type ComparisonResults []ComparisonResult

// GenerateAlterStatements generates alter statements from comparison results.
func (p ComparisonResults) GenerateAlterStatements(w io.Writer, sg StatementGenerator) {
	for _, result := range p {
		w.Write([]byte(result.GenerateStatement(sg)))
		w.Write([]byte{'\n'})
	}
	return
}

// ComparisonResult represents a comparison result.
type ComparisonResult interface {
	GenerateStatement(sg StatementGenerator) string
}

// TableMissingComparisonResult represents the comparison result that the table is missing.
type TableMissingComparisonResult struct {
	Table *Table
}

// GenerateStatement generates alter statement from the comparison result.
func (p *TableMissingComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateCreateTableStatement(p.Table)
}

// TableRedundantComparisonResult represents the comparison result that the table is redundant.
type TableRedundantComparisonResult struct {
	TableName string
}

// GenerateStatement generates alter statement from the comparison result.
func (p *TableRedundantComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateDropTableStatement(p.TableName)
}

// TableChangedComparisonResult represents the comparison result that the table is changed.
type TableChangedComparisonResult struct {
	Table    *Table
	OldTable *Table
}

// GenerateStatement generates alter statement from the comparison result.
func (p *TableChangedComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateDropTableStatement(p.OldTable.Name) + "\n" +
		sg.GenerateCreateTableStatement(p.Table)
}

// TableCommentChangedComparisonResult represents the comparison result that the table comment is changed.
type TableCommentChangedComparisonResult struct {
	Table *Table
}

// GenerateStatement generates alter statement from the comparison result.
func (p *TableCommentChangedComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateAlterTableCommentStatement(p.Table)
}

// ColumnMissingComparisonResult represents the comparison result that the column is missing.
type ColumnMissingComparisonResult struct {
	Table       *Table
	ColumnIndex int
}

// GenerateStatement generates alter statement from the comparison result.
func (p *ColumnMissingComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateAddColumnStatement(p.Table, p.ColumnIndex)
}

// ColumnRedundantComparisonResult represents the comparison result the column is redundant.
type ColumnRedundantComparisonResult struct {
	TableName  string
	ColumnName string
}

// GenerateStatement generates alter statement from the comparison result.
func (p *ColumnRedundantComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateDropColumnStatement(p.TableName, p.ColumnName)
}

// ColumnChangedComparisonResult represents the comparison result that the column is changed.
type ColumnChangedComparisonResult struct {
	Table       *Table
	ColumnIndex int
}

// GenerateStatement generates alter statement from the comparison result.
func (p *ColumnChangedComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateModifyColumnStatement(p.Table, p.ColumnIndex)
}

// PrimaryKeyMissingComparisonResult represents the comparison result that the primary key is missing.
type PrimaryKeyMissingComparisonResult struct {
	Table *Table
}

// GenerateStatement generates alter statement from the comparison result.
func (p *PrimaryKeyMissingComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateAddPrimaryKeyStatement(p.Table)
}

// PrimaryKeyRedundantComparisonResult represents that the primary key is redundant.
type PrimaryKeyRedundantComparisonResult struct {
	TableName string
}

// GenerateStatement generates alter statement from the comparison result.
func (p *PrimaryKeyRedundantComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateDropPrimaryKeyStatement(p.TableName)
}

// PrimaryKeyChangedComparisonResult represents the comparison result that the primary key is changed.
type PrimaryKeyChangedComparisonResult struct {
	Table *Table
}

// GenerateStatement generates alter statement from the comparison result.
func (p *PrimaryKeyChangedComparisonResult) GenerateStatement(sg StatementGenerator) string {
	return sg.GenerateDropPrimaryKeyStatement(p.Table.Name) + "\n" +
		sg.GenerateAddPrimaryKeyStatement(p.Table)
}

// Compare compares the table with its old version and generates Statement for altering database.
func Compare(table *Table, oldTable *Table) (results ComparisonResults) {
	results = make([]ComparisonResult, 0)
	if table == nil && oldTable == nil {
		return
	}
	if table != nil && oldTable == nil {
		results = append(results, &TableMissingComparisonResult{Table: table})
		return
	}
	if table == nil && oldTable != nil {
		results = append(results, &TableRedundantComparisonResult{TableName: oldTable.Name})
		return
	}
	if table.Name != oldTable.Name {
		results = append(results, &TableChangedComparisonResult{Table: table, OldTable: oldTable})
		return
	}
	if table.Comment != oldTable.Comment {
		results = append(results, &TableCommentChangedComparisonResult{Table: table})
	}
	for i, column := range table.Columns {
		oldColumn := findColumn(oldTable, column.Name)
		if oldColumn == nil {
			results = append(results, &ColumnMissingComparisonResult{Table: table, ColumnIndex: i})
		} else {
			if !columnEqual(column, oldColumn) {
				results = append(results, &ColumnChangedComparisonResult{Table: table, ColumnIndex: i})
			}
		}
	}
	for _, oldColumn := range oldTable.Columns {
		column := findColumn(table, oldColumn.Name)
		if column == nil {
			results = append(results, &ColumnRedundantComparisonResult{TableName: table.Name, ColumnName: oldColumn.Name})
		}
	}
	pkColumns := table.PrimaryKeyColumnNames()
	oldPKColumns := oldTable.PrimaryKeyColumnNames()
	if len(pkColumns) > 0 && len(oldPKColumns) == 0 {
		results = append(results, &PrimaryKeyMissingComparisonResult{Table: table})
	} else if len(pkColumns) == 0 && len(oldPKColumns) > 0 {
		results = append(results, &PrimaryKeyRedundantComparisonResult{TableName: table.Name})
	} else if len(pkColumns) > 0 && len(oldPKColumns) > 0 && !comm.StringArrayEqual(pkColumns, oldPKColumns) {
		results = append(results, &PrimaryKeyChangedComparisonResult{Table: table})
	}
	return
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
