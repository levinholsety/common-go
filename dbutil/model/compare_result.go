package model

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
