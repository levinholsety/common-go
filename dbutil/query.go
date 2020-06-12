package dbutil

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

var (
	errInvalidType          = errors.New("invalid type")
	errStructNotAppropriate = errors.New("struct should have 'tbl' and 'col' tags")
)

// NewQuery creates a query from specified struct type.
// The struct should have 'tbl' and 'col' tags.
func NewQuery(elemType reflect.Type) (query *Query, err error) {
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		err = errInvalidType
		return
	}
	var tblName string
	var colNames []string
	var colFieldIndexes []int
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if len(tblName) == 0 {
			tblName = field.Tag.Get("tbl")
		}
		colName := field.Tag.Get("col")
		if len(colName) > 0 {
			colNames = append(colNames, colName)
			colFieldIndexes = append(colFieldIndexes, i)
		}
	}
	if len(tblName) == 0 || len(colNames) == 0 {
		err = errStructNotAppropriate
		return
	}
	query = &Query{
		qry: "select " + strings.Join(colNames, ",") + " from " + tblName,
		rec: func(rows *sql.Rows) (rec interface{}, err error) {
			elem := reflect.New(elemType)
			values := make([]interface{}, len(colFieldIndexes))
			for i, colFieldIndex := range colFieldIndexes {
				values[i] = elem.Elem().Field(colFieldIndex).Addr().Interface()
			}
			err = rows.Scan(values...)
			if err != nil {
				return
			}
			rec = elem.Interface()
			return
		},
	}
	return
}

// Query provides simple method for quering database.
type Query struct {
	qry string
	rng *Range
	rec func(*sql.Rows) (interface{}, error)
}

// Append appends query string to current query.
func (p *Query) Append(queryString string) *Query {
	p.qry += " " + queryString
	return p
}

// SetRange sets query range to current query.
func (p *Query) SetRange(rng *Range) *Query {
	p.rng = rng
	return p
}

// Execute executes current query and invokes onRecord after a record is read.
// recIndex represents the index of current record. It starts with 0.
// rowIndex represents the index of row in all rows of current query. It starts with the offset of the query range.
func (p *Query) Execute(db *sql.DB, onRecord func(recIndex int, rowIndex int, rec interface{}), args ...interface{}) (err error) {
	rows, err := db.Query(p.qry, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	index := -1
	count := 0
	for rows.Next() {
		index++
		if p.rng != nil && index < p.rng.Offset {
			continue
		}
		var rec interface{}
		rec, err = p.rec(rows)
		if err != nil {
			return
		}
		onRecord(count, index, rec)
		count++
		if p.rng != nil && count == p.rng.Length {
			break
		}
	}
	return
}

// ExecuteSlice executes current query and returns the result as slice.
func (p *Query) ExecuteSlice(db *sql.DB, args ...interface{}) (result []interface{}, err error) {
	err = p.Execute(db, func(_, _ int, rec interface{}) {
		result = append(result, rec)
	}, args...)
	return
}
