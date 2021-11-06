package codegenerator

import (
	"fmt"
	"strings"
)

type TokenLiteral interface {
	Equals(other TokenLiteral) bool
}

type ColumnLiteral struct {
	Table  string
	Column string
}

func (c ColumnLiteral) Equals(other TokenLiteral) bool {
	if columnLiteral, ok := other.(ColumnLiteral); ok {
		return c.Table == columnLiteral.Table &&
			c.Column == columnLiteral.Column
	}

	return false
}

func (c *ColumnLiteral) String() string {
	return fmt.Sprintf(`"%s"."%s"`, c.Table, c.Column)
}

type FunctionLiteral struct {
	Name       string
	Type       TokenType
	Parameters []*Token
}

func (f FunctionLiteral) Equals(other TokenLiteral) bool {
	if functionLiteral, ok := other.(FunctionLiteral); ok {
		areParametersEqual := areTokenSlicesEqual(f.Parameters, functionLiteral.Parameters)
		return f.Type == functionLiteral.Type &&
			areParametersEqual
	}

	return false
}

func (f FunctionLiteral) String() string {
	paramString := ""
	for _, parameter := range f.Parameters {
		paramString += fmt.Sprintf("%s, ", parameter.Lexeme)
	}
	paramString = strings.Trim(paramString, ", ")

	return fmt.Sprintf("%s(%s)", f.Name, paramString)
}

type NumberLiteral float64

func (c NumberLiteral) Equals(other TokenLiteral) bool {
	if numberLiteral, ok := other.(NumberLiteral); ok {
		return c == numberLiteral
	}

	return false
}

func (c NumberLiteral) String() string {
	return fmt.Sprintf("%f", float64(c))
}

type StringLiteral string

func (s StringLiteral) Equals(other TokenLiteral) bool {
	if stringLiteral, ok := other.(StringLiteral); ok {
		return s == stringLiteral
	}

	return s == other
}

func (s StringLiteral) String() string {
	return string(s)
}
