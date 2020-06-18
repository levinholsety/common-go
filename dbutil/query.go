package dbutil

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

var (
	errInvalidType          = errors.New("invalid type")
	errStructNotAppropriate = errors.New("struct should have 'tbl' tag and at least one 'col' or 'exp' tag")
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
	var selectExpressions []string
	var fieldIndexes []int
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if len(tblName) == 0 {
			tblName = field.Tag.Get("tbl")
		}
		selectExpression := field.Tag.Get("exp")
		if len(selectExpression) > 0 {
			selectExpressions = append(selectExpressions, selectExpression)
			fieldIndexes = append(fieldIndexes, i)
		} else {
			selectExpression = field.Tag.Get("col")
			if len(selectExpression) > 0 {
				selectExpressions = append(selectExpressions, selectExpression)
				fieldIndexes = append(fieldIndexes, i)
			}
		}
	}
	if len(tblName) == 0 || len(selectExpressions) == 0 {
		err = errStructNotAppropriate
		return
	}
	query = &Query{
		queryString: "select " + strings.Join(selectExpressions, ",") + " from " + tblName,
		readRecord: func(rows *sql.Rows) (rec interface{}, err error) {
			elem := reflect.New(elemType)
			values := make([]interface{}, len(fieldIndexes))
			for i, fieldIndex := range fieldIndexes {
				values[i] = elem.Elem().Field(fieldIndex).Addr().Interface()
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
	queryString string
	queryRange  *Range
	readRecord  func(*sql.Rows) (interface{}, error)
}

// Append appends query string to current query.
func (p *Query) Append(queryString string) *Query {
	p.queryString += " " + queryString
	return p
}

// SetRange sets query range to current query.
func (p *Query) SetRange(rng *Range) *Query {
	p.queryRange = rng
	return p
}

// Execute executes current query and invokes onRecord after a record is read.
// recIndex represents the index of current record. It starts with 0.
// rowIndex represents the index of row in all rows of current query. It starts with the offset of the query range.
func (p *Query) Execute(db *sql.DB, onRecord func(recIndex int, rowIndex int, rec interface{}), args ...interface{}) (err error) {
	rows, err := db.Query(p.queryString, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	index := -1
	count := 0
	for rows.Next() {
		index++
		if p.queryRange != nil && index < p.queryRange.Offset {
			continue
		}
		var rec interface{}
		rec, err = p.readRecord(rows)
		if err != nil {
			return
		}
		onRecord(count, index, rec)
		count++
		if p.queryRange != nil && count == p.queryRange.Length {
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
