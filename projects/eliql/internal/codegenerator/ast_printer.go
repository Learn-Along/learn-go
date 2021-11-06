package codegenerator

import (
	"fmt"
	"strings"
)

// AstPrinter is a Visitor implementation to print Abstract Syntax Tree (AST) in lisp style
type AstPrinter struct {
}

func (a *AstPrinter) String(e Expression) Output {
	return e.accept(a)
}

func (a *AstPrinter) visitUnionExpression(u *UnionExpression) Output {
	exprs := make([]Expression, len(u.SelectExprs))
	for i, expr := range u.SelectExprs {
		exprs[i] = expr
	}

	return a.parenthesize("union", exprs...)
}

func (a *AstPrinter) visitSelectExpression(s *SelectExpression) Output {
	exprs := make([]Expression, 0)
	if s.ColumnExprs != nil {
		for _, expr := range s.ColumnExprs {
			exprs = append(exprs, expr)
		}
	}

	if s.JoinExprs != nil {
		for _, expr := range s.JoinExprs {
			exprs = append(exprs, expr)
		}
	}

	if s.WhereExpr != nil {
		exprs = append(exprs, s.WhereExpr)
	}

	if s.GroupByExpr != nil {
		exprs = append(exprs, s.GroupByExpr)
	}

	if s.OrderByExpr != nil {
		exprs = append(exprs, s.OrderByExpr)
	}

	return a.parenthesize(fmt.Sprintf("select from %s", s.Table.Lexeme), exprs...)
}

func (a *AstPrinter) visitUnionSelectExpression(u *UnionSelectExpression) Output {
	all := ""
	if u.All != nil {
		all = "ALL"
	}
	
	return a.parenthesize(all, u.SelectExpr)
}

func (a *AstPrinter) visitColumnExpression(c *ColumnExpression) Output {
	exprs := make([]Expression, 0)
	if c.ArithmeticExpr != nil {
		exprs = append(exprs, c.ArithmeticExpr)
	}

	title := ""
	if c.Name != nil {
		title = fmt.Sprintf("define %s", c.Name.Literal.(StringLiteral))
	}

	if c.Column != nil {
		title = fmt.Sprintf("%s %s", title, c.Column.Literal.(ColumnLiteral))
	}

	if c.Function != nil {
		title = fmt.Sprintf("%s %s", title, c.Function.Literal.(FunctionLiteral))
	}

	return a.parenthesize(strings.Trim(title, " "), exprs...)
}

func (a *AstPrinter) visitJoinExpression(j *JoinExpression) Output {
	exprs := make([]Expression, 0)
	if j.Conditions != nil {
		for _, condition := range j.Conditions {
			exprs = append(exprs, condition)
		}
	}

	title := fmt.Sprintf("%s JOIN %s", j.Type.Lexeme, j.Table.Literal)
	return a.parenthesize(title, exprs...)
}

func (a *AstPrinter) visitWhereExpression(w *WhereExpression) Output {
	exprs := make([]Expression, 0)
	if w.Comparisons != nil {
		for _, comparison := range w.Comparisons {
			exprs = append(exprs, comparison)
		}
	}

	return a.parenthesize("WHERE", exprs...)
}

func (a *AstPrinter) visitGroupByExpression(g *GroupByExpression) Output {
	representation := "GROUP BY "
	for _, column := range g.Columns {
		representation += fmt.Sprintf("%s ", column.Literal.(ColumnLiteral))
	}

	return strings.Trim(representation, " ")
}

func (a *AstPrinter) visitOrderByExpression(o *OrderByExpression) Output {
	exprs := make([]Expression, 0)
	if o.ColumnOrderExprs != nil {
		for _, expr := range o.ColumnOrderExprs {
			exprs = append(exprs, expr)
		}
	}

	return a.parenthesize("ORDER BY", exprs...)
}

func (a *AstPrinter) visitJoinCondition(j *JoinCondition) Output {
	return fmt.Sprintf("(= %s %s)", j.Left.Literal.(ColumnLiteral), j.Right.Literal.(ColumnLiteral))
}

func (a *AstPrinter) visitComparisonExpression(c *ComparisonExpression) Output {
	logicalOperator := "AND"
	if c.LogicalOperator != nil {
		logicalOperator = c.LogicalOperator.Lexeme
	}

	comparison := a.parenthesize(fmt.Sprintf("%s %s", c.Comparator.Lexeme, c.Left.Literal.(ColumnLiteral)), c.Right)
	if c.Not != nil {
		comparison = fmt.Sprintf("(NOT %s)", comparison)
	}

	return fmt.Sprintf("(%s %s)", logicalOperator, comparison)
}

func (a *AstPrinter) visitColumnOrderExpression(c *ColumnOrderExpression) Output {
	return fmt.Sprintf("(%s %s)", c.Order.Lexeme, c.Column.Literal.(ColumnLiteral))
}

func (a *AstPrinter) visitArithmeticExpression(ar *ArithmeticExpression) Output {
	return fmt.Sprintf("(%s %s %s)", ar.Operator.Lexeme, ar.Left.Literal, ar.Right.Literal)
}

func (a *AstPrinter) parenthesize(s string, exprs ...Expression) Output {
	tree := ""
	for _, expr := range exprs {
		tree += fmt.Sprintf("%v", expr.accept(a))
	}

	return fmt.Sprintf("(%s, %s)", s, tree)
}
