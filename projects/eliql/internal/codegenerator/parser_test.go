package codegenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParser(t *testing.T) {
	tokens := []Token{
		{
			Type:   String,
			Lexeme:  "Lexeme",
			Literal: StringLiteral("literal"),
			Line:    1,
		},		
		{
			Type:   Number,
			Lexeme:  "Lexeme",
			Literal: NumberLiteral(9.6),
			Line:    156,
		},
		{
			Type:   MaxFunc,
			Lexeme:  `MAX("table"."column")`,
			Literal: FunctionLiteral{
				Name:       "MAX",
				Type:       MaxFunc,
				Parameters: []*Token{},
			},
			Line:    200,
		},
	}

	parser := NewParser(tokens)
	assert.Equal(t, parser.tokens, tokens)
	assert.Equal(t, parser.current, 0)
}

func TestParser_expression(t *testing.T) {
	// SELECT "foo"."bar" AS 'bar' FROM "foo";
	tokens := []Token{
		{
			Type:   Select,
			Lexeme:  "SELECT",
			Line:    1,
		},
		{
			Type:   Column,
			Lexeme:  `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table: "foo",
				Column: "bar",
			},
			Line:    1,
		},
		{
			Type:    As,
			Lexeme:  "AS",
			Line:    1,
		},
		{
			Type:    String,
			Lexeme:  "'bar'",
			Literal: StringLiteral("bar"),
			Line:    1,
		},
		{
			Type:   From,
			Lexeme:  "FROM",
			Line:    1,
		},
		{
			Type:   Table,
			Lexeme:  "foo",
			Line:    1,
		},
		{
			Type:   SemiColon,
			Lexeme:  ";",
			Line:    1,
		},
	}

	parser := Parser{
		tokens:  tokens,
		current: 0,
	}

	expectedExpr := &SelectExpression{
		ColumnExprs: []*ColumnExpression{
			{
				Column:         &Token{
					Type:    Column,
					Lexeme:  `"foo"."bar"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "bar",
					},
					Line:    1,
				},
				Name:           &Token{
					Type:    String,
					Lexeme:  "'bar'",
					Literal: StringLiteral("bar"),
					Line:    1,
				},
			},
		},
		Table:       &Token{
			Type:    Table,
			Lexeme:  "foo",
			Line:    1,
		},
	}

	actualExpr := parser.expression().(*SelectExpression)
	assert.True(t, areSelectExpressionsEqual(expectedExpr, actualExpr))
}

func TestParser_selectExpr(t *testing.T) {
	// SELECT "foo"."bar" AS 'bar' FROM "foo";
	tokens := []Token{
		{
			Type:   Select,
			Lexeme:  "SELECT",
			Line:    1,
		},
		{
			Type:   Column,
			Lexeme:  `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table: "foo",
				Column: "bar",
			},
			Line:    1,
		},
		{
			Type:    As,
			Lexeme:  "AS",
			Line:    1,
		},
		{
			Type:    String,
			Lexeme:  "'bar'",
			Literal: StringLiteral("bar"),
			Line:    1,
		},
		{
			Type:   From,
			Lexeme:  "FROM",
			Line:    1,
		},
		{
			Type:   Table,
			Lexeme:  "foo",
			Line:    1,
		},
		{
			Type:   SemiColon,
			Lexeme:  ";",
			Line:    1,
		},
	}

	parser := Parser{
		tokens:  tokens,
		current: 0,
	}

	expectedExpr := &SelectExpression{
		ColumnExprs: []*ColumnExpression{
			{
				Column:         &Token{
					Type:    Column,
					Lexeme:  `"foo"."bar"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "bar",
					},
					Line:    1,
				},
				Name:           &Token{
					Type:    String,
					Lexeme:  "'bar'",
					Literal: StringLiteral("bar"),
					Line:    1,
				},
			},
		},
		Table:       &Token{
			Type:    Table,
			Lexeme:  "foo",
			Line:    1,
		},
	}

	actualExpr := parser.selectExpr()
	assert.True(t, areSelectExpressionsEqual(expectedExpr, actualExpr))
}

func TestParser_columnExpr(t *testing.T) {
	// "foo"."bar" AS 'bar'
	tokens := []Token{
		{
			Type:   Column,
			Lexeme:  `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table: "foo",
				Column: "bar",
			},
			Line:    1,
		},
		{
			Type:    As,
			Lexeme:  "AS",
			Line:    1,
		},
		{
			Type:    String,
			Lexeme:  "'bar'",
			Literal: StringLiteral("bar"),
			Line:    1,
		},
	}

	parser := Parser{
		tokens:  tokens,
		current: 0,
	}

	expectedExpr := &ColumnExpression{
				Column:         &Token{
					Type:    Column,
					Lexeme:  `"foo"."bar"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "bar",
					},
					Line:    1,
				},
				Name:           &Token{
					Type:    String,
					Lexeme:  "'bar'",
					Literal: StringLiteral("bar"),
					Line:    1,
				},
			}

	actualExpr := parser.columnExpr()
	assert.True(t, areColumnExpressionsEqual(expectedExpr, actualExpr))
}

func areSelectExpressionsEqual(expr1 *SelectExpression, expr2 *SelectExpression) bool {
	if *(expr1.Table) != *(expr2.Table) {
		return false
	}

	if len(expr1.ColumnExprs) != len(expr2.ColumnExprs) {
		return false
	}

	for i := range expr1.ColumnExprs {
		if !areColumnExpressionsEqual(expr1.ColumnExprs[i], expr2.ColumnExprs[i]) {
			return false
		}
	}
	return true
}

func areColumnExpressionsEqual(expr1 *ColumnExpression, expr2 *ColumnExpression) bool {
	return *(expr1.Name) == *(expr2.Name) &&
		*(expr1.Column) == *(expr2.Column)
}