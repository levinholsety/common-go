package mysql

import (
	"bytes"
	"fmt"

	"github.com/levinholsety/common-go/dbutil/model"
	"github.com/levinholsety/common-go/utils"
)

// StatementGenerator represents the statement generator for MySQL.
type StatementGenerator struct{}

var _ model.StatementGenerator = new(StatementGenerator)

// GenerateUseStatement generates use database Statement.
func (p *StatementGenerator) GenerateUseStatement(schema *model.Schema) string {
	return fmt.Sprintf("USE `%s`;", schema.Name)
}

func (p *StatementGenerator) columnStatement(column *model.Column) (result string) {
	result = fmt.Sprintf("`%s` %s", column.Name, column.Type)
	if column.Nullable {
		result += " NULL"
	} else {
		result += " NOT NULL"
	}
	if len(column.Comment) > 0 {
		result += " COMMENT '" + column.Comment + "'"
	}
	return
}

func (p *StatementGenerator) primaryKeyStatement(table *model.Table) (result string) {
	for _, column := range table.Columns {
		if column.IsPrimaryKey {
			if len(result) > 0 {
				result += ","
			}
			result += "`" + column.Name + "`"
		}
	}
	if len(result) > 0 {
		result = "PRIMARY KEY (" + result + ")"
	}
	return
}

// GenerateCreateTableStatement generates create table Statement.
func (p *StatementGenerator) GenerateCreateTableStatement(table *model.Table) string {
	buf := &bytes.Buffer{}
	w := utils.NewTextWriter(buf)
	w.WriteLineFormat("CREATE TABLE `%s` (", table.Name)
	for i, column := range table.Columns {
		if i > 0 {
			w.WriteLine(",")
		}
		w.WriteString("    ")
		w.WriteString(p.columnStatement(column))
	}
	pkStr := p.primaryKeyStatement(table)
	if len(pkStr) > 0 {
		w.WriteLine(",")
		w.WriteString("    ")
		w.WriteLine(pkStr)
	}
	w.WriteString(")")
	if len(table.Comment) > 0 {
		w.WriteString(" COMMENT='")
		w.WriteString(table.Comment)
		w.WriteString("'")
	}
	w.WriteString(";")
	return buf.String()
}

// GenerateAlterTableCommentStatement generates alter table comment statement.
func (p *StatementGenerator) GenerateAlterTableCommentStatement(table *model.Table) string {
	return fmt.Sprintf("ALTER TABLE `%s` COMMENT='%s';", table.Name, table.Comment)
}

// GenerateDropTableStatement generates drop table statement.
func (p *StatementGenerator) GenerateDropTableStatement(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", tableName)
}

// GenerateAddColumnStatement generates add column statement.
func (p *StatementGenerator) GenerateAddColumnStatement(table *model.Table, columnIndex int) string {
	column := table.Columns[columnIndex]
	if columnIndex == 0 {
		return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s FIRST;", table.Name, p.columnStatement(column))
	}
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s AFTER `%s`;", table.Name, p.columnStatement(column), table.Columns[columnIndex-1].Name)
}

// GenerateModifyColumnStatement generates modify column statement.
func (p *StatementGenerator) GenerateModifyColumnStatement(table *model.Table, columnIndex int) string {
	column := table.Columns[columnIndex]
	if columnIndex == 0 {
		return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s FIRST;", table.Name, p.columnStatement(column))
	}
	return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s AFTER `%s`;", table.Name, p.columnStatement(column), table.Columns[columnIndex-1].Name)
}

// GenerateDropColumnStatement generates drop column statement.
func (p *StatementGenerator) GenerateDropColumnStatement(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", tableName, columnName)
}

// GenerateAddPrimaryKeyStatement generates add primary key statement.
func (p *StatementGenerator) GenerateAddPrimaryKeyStatement(table *model.Table) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD %s;", table.Name, p.primaryKeyStatement(table))
}

// GenerateDropPrimaryKeyStatement generates drop primary key statement.
func (p *StatementGenerator) GenerateDropPrimaryKeyStatement(tableName string) string {
	return fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY;", tableName)
}
