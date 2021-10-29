package codegenerator

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

type FunctionLiteral struct {
	Type       TokenType
	Parameters []*Token
}
func (c FunctionLiteral) Equals(other TokenLiteral) bool {
	if functionLiteral, ok := other.(FunctionLiteral); ok {
		areParametersEqual := areTokenSlicesEqual(c.Parameters, functionLiteral.Parameters)
		return c.Type == functionLiteral.Type &&
			areParametersEqual
	}

	return false
}

type NumberLiteral float64
func (c NumberLiteral) Equals(other TokenLiteral) bool {
	if numberLiteral, ok := other.(NumberLiteral); ok {
		return c == numberLiteral
	}

	return false
}

type StringLiteral string
func (s StringLiteral) Equals(other TokenLiteral) bool {
	if stringLiteral, ok := other.(StringLiteral); ok {
		return s == stringLiteral
	}

	return s == other
}
