package model

// DataClass represents the class of data type.
type DataClass int

// DataClasses.
const (
	Text DataClass = iota + 1
	Number
	Binary
	Time
)

// Model represents database model.
type Model struct {
	Schemas []*Schema `json:"schemas"`
}

// Schema represents database schema.
type Schema struct {
	Name   string   `json:"name"`
	Tables []*Table `json:"tables"`
}

// Table represents database table.
type Table struct {
	Name    string    `json:"name"`
	Comment string    `json:"comment,omitempty"`
	Columns []*Column `json:"columns"`
}

// PrimaryKeyColumnNames returns the names of primary key columns.
func (p *Table) PrimaryKeyColumnNames() (result []string) {
	for _, column := range p.Columns {
		if column.IsPrimaryKey {
			result = append(result, column.Name)
		}
	}
	return
}

// Column represents database table column.
type Column struct {
	Name         string    `json:"name"`
	DataType     string    `json:"dataType"`
	DataClass    DataClass `json:"dataClass"`
	Type         string    `json:"type"`
	Nullable     bool      `json:"nullable"`
	IsPrimaryKey bool      `json:"isPrimaryKey"`
	Comment      string    `json:"comment,omitempty"`
}
