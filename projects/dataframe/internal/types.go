package internal

import (
	"regexp"

	"github.com/tobgu/qframe"
)

// type FilterType struct {
// 	clause qframe.FilterClause
// }

type View interface {
	ItemAt(i int) bool
	Len() int
	Slice() []bool
}

type Column interface {
	GreaterThan(operand float64) qframe.FilterClause
	GreaterOrEquals(operand float64) qframe.FilterClause
	LessThan(operand float64) qframe.FilterClause
	LessOrEquals(operand float64) qframe.FilterClause
	Equals(operand interface{}) qframe.FilterClause
	IsLike(pattern *regexp.Regexp) qframe.FilterClause
}

type ColumnComparator struct {
	Name string
}
func (c ColumnComparator) GreaterThan(operand float64) qframe.FilterClause {
	return qframe.Filter{Column: c.Name, Comparator: ">", Arg: operand}
}
func (c ColumnComparator) GreaterOrEquals(operand float64) qframe.FilterClause {
	return qframe.Filter{Column: c.Name, Comparator: ">=", Arg: operand}
}
func (c ColumnComparator) LessThan(operand float64) qframe.FilterClause {
	return qframe.Filter{Column: c.Name, Comparator: "<", Arg: operand}
}
func (c ColumnComparator) LessOrEquals(operand float64) qframe.FilterClause {
	return qframe.Filter{Column: c.Name, Comparator: "<=", Arg: operand}
}
func (c ColumnComparator) Equals(operand interface{}) qframe.FilterClause {
	return qframe.Filter{Column: c.Name, Comparator: "=", Arg: operand}
}
func (c ColumnComparator) IsLike(pattern *regexp.Regexp) qframe.FilterClause {
	matches := func(x string) bool { return pattern.MatchString(x) }

	return qframe.Filter{Column: c.Name, Comparator: matches, Arg: pattern}
}

type IntColumn struct {
	qframe.IntView
	ColumnComparator
}

// Creates a new IntColumn
func NewIntColumn(q qframe.QFrame, name string) (*IntColumn, error) {
	v, err := q.IntView(name)
	if err != nil {
		return nil, err
	}

	return &IntColumn{v, ColumnComparator{Name: name}}, nil
}

type FloatColumn struct {
	qframe.FloatView
	ColumnComparator
}

// Creates a new IntColumn
func NewFloatColumn(q qframe.QFrame, name string) (*FloatColumn, error) {
	v, err := q.FloatView(name)
	if err != nil {
		return nil, err
	}

	return &FloatColumn{v, ColumnComparator{Name: name}}, nil
}

type StringColumn struct {
	qframe.StringView
	ColumnComparator
}

// Creates a new IntColumn
func NewStringColumn(q qframe.QFrame, name string) (*StringColumn, error) {
	v, err := q.StringView(name)
	if err != nil {
		return nil, err
	}

	return &StringColumn{v, ColumnComparator{Name: name}}, nil
}

type BoolColumn struct {
	qframe.BoolView
	ColumnComparator
}

// Creates a new IntColumn
func NewBoolColumn(q qframe.QFrame, name string) (*BoolColumn, error) {
	v, err := q.BoolView(name)
	if err != nil {
		return nil, err
	}

	return &BoolColumn{v, ColumnComparator{Name: name}}, nil
}