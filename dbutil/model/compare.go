package model

import "bytes"

// ComparisonResults represents comparison results.
type ComparisonResults []ComparisonResult

// GenerateAlterStatements generates alter statements from comparison results.
func (p ComparisonResults) GenerateAlterStatements(sg StatementGenerator) string {
	buf := &bytes.Buffer{}
	for _, result := range p {
		buf.WriteString(result.GenerateStatement(sg))
		buf.WriteByte('\n')
	}
	return buf.String()
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
		oldColumn, ok := findColumn(oldTable, column.Name)
		if !ok {
			results = append(results, &ColumnMissingComparisonResult{Table: table, ColumnIndex: i})
		} else {
			if !isColumnSame(column, oldColumn) {
				results = append(results, &ColumnChangedComparisonResult{Table: table, ColumnIndex: i})
			}
		}
	}
	for _, oldColumn := range oldTable.Columns {
		_, ok := findColumn(table, oldColumn.Name)
		if !ok {
			results = append(results, &ColumnRedundantComparisonResult{TableName: table.Name, ColumnName: oldColumn.Name})
		}
	}
	pkColumns := table.PrimaryKeyColumns()
	oldPKColumns := oldTable.PrimaryKeyColumns()
	if len(pkColumns) > 0 && len(oldPKColumns) == 0 {
		results = append(results, &PrimaryKeyMissingComparisonResult{Table: table})
	} else if len(pkColumns) == 0 && len(oldPKColumns) > 0 {
		results = append(results, &PrimaryKeyRedundantComparisonResult{TableName: table.Name})
	} else if len(pkColumns) > 0 && len(oldPKColumns) > 0 && !isColumnsSame(pkColumns, oldPKColumns) {
		results = append(results, &PrimaryKeyChangedComparisonResult{Table: table})
	}
	return
}

func findColumn(table *Table, columnName string) (column *Column, ok bool) {
	for _, column = range table.Columns {
		if column.Name == columnName {
			ok = true
			return
		}
	}
	ok = false
	return
}

func isColumnSame(column1, column2 *Column) bool {
	return column1.Name == column2.Name &&
		column1.Type == column2.Type &&
		column1.Nullable == column2.Nullable &&
		column1.Comment == column2.Comment
}

func isColumnsSame(columns1, columns2 []*Column) bool {
	if len(columns1) != len(columns2) {
		return false
	}
	for i, col1 := range columns1 {
		col2 := columns2[i]
		if col1.Name != col2.Name {
			return false
		}
	}
	return true
}
