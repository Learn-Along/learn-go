package codegenerator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParser(t *testing.T) {
	tokens := []Token{
		{
			Type:    String,
			Lexeme:  "Lexeme",
			Literal: StringLiteral("literal"),
			Line:    1,
		},
		{
			Type:    Number,
			Lexeme:  "Lexeme",
			Literal: NumberLiteral(9.6),
			Line:    156,
		},
		{
			Type:   MaxFunc,
			Lexeme: `MAX("table"."column")`,
			Literal: FunctionLiteral{
				Name:       "MAX",
				Type:       MaxFunc,
				Parameters: []*Token{},
			},
			Line: 200,
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
			Lexeme: "SELECT",
			Line:   1,
		},
		{
			Type:   Column,
			Lexeme: `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "bar",
			},
			Line: 1,
		},
		{
			Type:   As,
			Lexeme: "AS",
			Line:   1,
		},
		{
			Type:    String,
			Lexeme:  "'bar'",
			Literal: StringLiteral("bar"),
			Line:    1,
		},
		{
			Type:   From,
			Lexeme: "FROM",
			Line:   1,
		},
		{
			Type:   Table,
			Lexeme: "foo",
			Line:   1,
		},
		{
			Type:   SemiColon,
			Lexeme: ";",
			Line:   1,
		},
	}

	parser := Parser{
		tokens:  tokens,
		current: 0,
	}

	expectedExpr := &SelectExpression{
		ColumnExprs: []*ColumnExpression{
			{
				Column: &Token{
					Type:   Column,
					Lexeme: `"foo"."bar"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "bar",
					},
					Line: 1,
				},
				Name: &Token{
					Type:    String,
					Lexeme:  "'bar'",
					Literal: StringLiteral("bar"),
					Line:    1,
				},
			},
		},
		Table: &Token{
			Type:   Table,
			Lexeme: "foo",
			Line:   1,
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
			Lexeme: "SELECT",
			Line:   1,
		},
		{
			Type:   Column,
			Lexeme: `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "bar",
			},
			Line: 1,
		},
		{
			Type:   As,
			Lexeme: "AS",
			Line:   1,
		},
		{
			Type:    String,
			Lexeme:  "'bar'",
			Literal: StringLiteral("bar"),
			Line:    1,
		},
		{
			Type:   From,
			Lexeme: "FROM",
			Line:   1,
		},
		{
			Type:   Table,
			Lexeme: "foo",
			Line:   1,
		},
		{
			Type:   SemiColon,
			Lexeme: ";",
			Line:   1,
		},
	}

	parser := Parser{
		tokens:  tokens,
		current: 0,
	}

	expectedExpr := &SelectExpression{
		ColumnExprs: []*ColumnExpression{
			{
				Column: &Token{
					Type:   Column,
					Lexeme: `"foo"."bar"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "bar",
					},
					Line: 1,
				},
				Name: &Token{
					Type:    String,
					Lexeme:  "'bar'",
					Literal: StringLiteral("bar"),
					Line:    1,
				},
			},
		},
		Table: &Token{
			Type:   Table,
			Lexeme: "foo",
			Line:   1,
		},
	}

	actualExpr := parser.selectExpr()
	assert.True(t, areSelectExpressionsEqual(expectedExpr, actualExpr))
}

func TestParser_unionExpr(t *testing.T) {
	// SELECT "foo"."bar" AS 'bar' FROM "foo"
	// UNION
	// SELECT "pumpkin"."color" AS 'color' FROM "pumpkin"
	// UNION ALL
	// SELECT "pumpkin"."origin" AS 'orn' FROM "pumpkin";
	tokens := []Token{
		{
			Type:   Select,
			Lexeme: "SELECT",
			Line:   1,
		},
		{
			Type:   Column,
			Lexeme: `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "bar",
			},
			Line: 1,
		},
		{
			Type:   As,
			Lexeme: "AS",
			Line:   1,
		},
		{
			Type:    String,
			Lexeme:  "'bar'",
			Literal: StringLiteral("bar"),
			Line:    1,
		},
		{
			Type:   From,
			Lexeme: "FROM",
			Line:   1,
		},
		{
			Type:   Table,
			Lexeme: "foo",
			Line:   1,
		},
		{
			Type:   Union,
			Lexeme: "UNION",
			Line:   2,
		},
		{
			Type:   Select,
			Lexeme: "SELECT",
			Line:   3,
		},
		{
			Type:   Column,
			Lexeme: `"pumpkin"."color"`,
			Literal: ColumnLiteral{
				Table:  "pumpkin",
				Column: "color",
			},
			Line: 3,
		},
		{
			Type:   As,
			Lexeme: "AS",
			Line:   3,
		},
		{
			Type:    String,
			Lexeme:  "'color'",
			Literal: StringLiteral("color"),
			Line:    3,
		},
		{
			Type:   From,
			Lexeme: "FROM",
			Line:   3,
		},
		{
			Type:   Table,
			Lexeme: "pumpkin",
			Line:   3,
		},
		{
			Type:   Union,
			Lexeme: "UNION",
			Line:   4,
		},
		{
			Type:   All,
			Lexeme: "ALL",
			Line:   4,
		},
		{
			Type:   Select,
			Lexeme: "SELECT",
			Line:   5,
		},
		{
			Type:   Column,
			Lexeme: `"pumpkin"."origin"`,
			Literal: ColumnLiteral{
				Table:  "pumpkin",
				Column: "origin",
			},
			Line: 5,
		},
		{
			Type:   As,
			Lexeme: "AS",
			Line:   5,
		},
		{
			Type:    String,
			Lexeme:  "'orn'",
			Literal: StringLiteral("orn"),
			Line:    5,
		},
		{
			Type:   From,
			Lexeme: "FROM",
			Line:   5,
		},
		{
			Type:   Table,
			Lexeme: "pumpkin",
			Line:   5,
		},
		{
			Type:   SemiColon,
			Lexeme: ";",
			Line:   5,
		},
	}

	parser := Parser{
		tokens:  tokens,
		current: 6,
	}

	left := &SelectExpression{
		ColumnExprs: []*ColumnExpression{
			{
				Column: &Token{
					Type:   Column,
					Lexeme: `"foo"."bar"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "bar",
					},
					Line: 1,
				},
				Name: &Token{
					Type:    String,
					Lexeme:  "'bar'",
					Literal: StringLiteral("bar"),
					Line:    1,
				},
			},
		},
		Table: &Token{
			Type:   Table,
			Lexeme: "foo",
			Line:   1,
		},
	}

	middle := &SelectExpression{
		ColumnExprs: []*ColumnExpression{
			{
				Column: &Token{
					Type:   Column,
					Lexeme: `"pumpkin"."color"`,
					Literal: ColumnLiteral{
						Table:  "pumpkin",
						Column: "color",
					},
					Line: 3,
				},
				Name: &Token{
					Type:    String,
					Lexeme:  "'color'",
					Literal: StringLiteral("color"),
					Line:    3,
				},
			},
		},
		Table: &Token{
			Type:   Table,
			Lexeme: "pumpkin",
			Line:   3,
		},
	}

	right := &SelectExpression{
		ColumnExprs: []*ColumnExpression{
			{
				Column: &Token{
					Type:   Column,
					Lexeme: `"pumpkin"."origin"`,
					Literal: ColumnLiteral{
						Table:  "pumpkin",
						Column: "origin",
					},
					Line: 5,
				},
				Name: &Token{
					Type:    String,
					Lexeme:  "'orn'",
					Literal: StringLiteral("orn"),
					Line:    5,
				},
			},
		},
		Table: &Token{
			Type:   Table,
			Lexeme: "pumpkin",
			Line:   5,
		},
	}

	expectedExpr := &UnionExpression{SelectExprs: []*UnionSelectExpression{
		{SelectExpr: left},
		{SelectExpr: middle},
		{SelectExpr: right, All: &Token{
			Type:   All,
			Lexeme: "ALL",
			Line:   4,
		}},
	}}

	actualExpr := parser.unionExpr(left)
	assert.True(t, areUnionExpressionsEqual(expectedExpr, actualExpr))
}

func TestParser_columnExpr(t *testing.T) {
	// "foo"."bar" AS 'bar'
	tokens := []Token{
		{
			Type:   Column,
			Lexeme: `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "bar",
			},
			Line: 1,
		},
		{
			Type:   As,
			Lexeme: "AS",
			Line:   1,
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
		Column: &Token{
			Type:   Column,
			Lexeme: `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "bar",
			},
			Line: 1,
		},
		Name: &Token{
			Type:    String,
			Lexeme:  "'bar'",
			Literal: StringLiteral("bar"),
			Line:    1,
		},
	}

	actualExpr := parser.columnExpr()
	assert.True(t, areColumnExpressionsEqual(expectedExpr, actualExpr))
}

func TestParser_joinExpr(t *testing.T) {
	// INNER JOIN "bar" ON "foo"."datetime" = "bar"."datetime" AND "foo"."name" = "bar"."fullName";
	tokens := []Token{
		{
			Type:   Inner,
			Lexeme: "INNER",
			Line:   1,
		},
		{
			Type:   Join,
			Lexeme: "JOIN",
			Line:   1,
		},
		{
			Type:   Table,
			Lexeme: `"bar"`,
			Line:   1,
		},
		{
			Type:   On,
			Lexeme: "ON",
			Line:   1,
		},
		{
			Type:   Column,
			Lexeme: `"foo"."datetime"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "datetime",
			},
			Line: 1,
		},
		{
			Type:   Equal,
			Lexeme: "=",
			Line:   1,
		},
		{
			Type:   Column,
			Lexeme: `"bar"."datetime"`,
			Literal: ColumnLiteral{
				Table:  "bar",
				Column: "datetime",
			},
			Line: 1,
		},
		{
			Type:   And,
			Lexeme: "AND",
			Line:   1,
		},
		{
			Type:   Column,
			Lexeme: `"foo"."name"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "name",
			},
			Line: 1,
		},
		{
			Type:   Equal,
			Lexeme: "=",
			Line:   1,
		},
		{
			Type:   Column,
			Lexeme: `"bar"."fullName"`,
			Literal: ColumnLiteral{
				Table:  "bar",
				Column: "fullName",
			},
			Line: 1,
		},
		{
			Type:   SemiColon,
			Lexeme: ";",
			Line:   1,
		},
	}
	parser := Parser{
		tokens:  tokens,
		current: 2,
	}

	expectedExpr := &JoinExpression{
		Type: &Token{
			Type:   Inner,
			Lexeme: "INNER",
			Line:   1,
		},
		Table: &Token{
			Type:   Table,
			Lexeme: `"bar"`,
			Line:   1,
		},
		Conditions: []*JoinCondition{
			{
				Left: &Token{
					Type:   Column,
					Lexeme: `"foo"."datetime"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "datetime",
					},
					Line: 1,
				},
				Right: &Token{
					Type:   Column,
					Lexeme: `"bar"."datetime"`,
					Literal: ColumnLiteral{
						Table:  "bar",
						Column: "datetime",
					},
					Line: 1,
				},
			},
			{
				Left: &Token{
					Type:   Column,
					Lexeme: `"foo"."name"`,
					Literal: ColumnLiteral{
						Table:  "foo",
						Column: "name",
					},
					Line: 1,
				},
				Right: &Token{
					Type:   Column,
					Lexeme: `"bar"."fullName"`,
					Literal: ColumnLiteral{
						Table:  "bar",
						Column: "fullName",
					},
					Line: 1,
				},
			},
		},
	}

	actualExpr := parser.joinExpr()
	assert.True(t, areJoinExpressionsEqual(expectedExpr, actualExpr))
}

func TestParser_whereExpr(t *testing.T) {
	// TODO
}

func TestParser_groupByExpr(t *testing.T) {
	// TODO
}

func TestParser_orderByExpr(t *testing.T) {
	// TODO
}

func TestParser_arithmeticExpr(t *testing.T) {
	// TODO
}

func TestParser_joinConditionExpr(t *testing.T) {
	// TODO
}

func TestParser_comparisonExpr(t *testing.T) {
	// TODO
}

func TestParser_columnOrderExpr(t *testing.T) {
	// TODO
}

func TestParser_primaryExpr(t *testing.T) {
	// TODO
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

func areUnionExpressionsEqual(expr1 *UnionExpression, expr2 *UnionExpression) bool {
	if len(expr1.SelectExprs) != len(expr2.SelectExprs) {
		return false
	}

	for i := range expr1.SelectExprs {
		if !areSelectExpressionsEqual(expr1.SelectExprs[i].SelectExpr, expr2.SelectExprs[i].SelectExpr) {
			return false
		}

		if !areTokenPtrEqual(expr1.SelectExprs[i].All, expr2.SelectExprs[i].All) {
			return false
		}
	}

	return true
}

func areColumnExpressionsEqual(expr1 *ColumnExpression, expr2 *ColumnExpression) bool {
	return *(expr1.Name) == *(expr2.Name) &&
		*(expr1.Column) == *(expr2.Column)
}

func areJoinExpressionsEqual(expr1 *JoinExpression, expr2 *JoinExpression) bool {
	if len(expr1.Conditions) != len(expr2.Conditions) {
		return false
	}

	for i := range expr1.Conditions {
		if !areJoinConditionsEqual(expr1.Conditions[i], expr2.Conditions[i]) {
			return false
		}
	}

	return areTokenPtrEqual(expr1.Table, expr2.Table) && areTokenPtrEqual(expr1.Type, expr2.Type)
}

func areJoinConditionsEqual(cond1 *JoinCondition, cond2 *JoinCondition) bool {
	return areTokenPtrEqual(cond1.Left, cond2.Left) && areTokenPtrEqual(cond1.Right, cond2.Right)
}
