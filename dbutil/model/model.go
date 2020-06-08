package model

// DataClass represents the class of data type.
type DataClass int

const (
	Text DataClass = iota + 1
	Number
	Binary
	Time
)

type Model struct {
	Schemas []*Schema
}

type Schema struct {
	Name   string
	Tables []*Table
}

type Table struct {
	Name    string
	Comment string
	Columns []*Column
}

func (p *Table) PrimaryKeyColumns() (result []*Column) {
	for _, column := range p.Columns {
		if column.IsPrimaryKey {
			result = append(result, column)
		}
	}
	return
}

type Column struct {
	Name         string
	DataType     string
	DataClass    DataClass
	Type         string
	Nullable     bool
	IsPrimaryKey bool
	Comment      string
}
