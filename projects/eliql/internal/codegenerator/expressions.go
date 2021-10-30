package codegenerator

type Expression interface {
	accept(visitor Visitor) Output
}

type Visitor interface {
	visitUnionExpr(u UnionExpression) Output
	visitSelectExpr(s SelectExpression) Output
}

type Output interface {}

type UnionExpression struct {
	SelectExprs []*UnionSelectExpression
}

type SelectExpression struct {
	ColumnExprs []*ColumnExpression
	Table *Token
	JoinExprs []*JoinExpression
	WhereExprs *WhereExpression
	GroupByExprs *GroupByExpression
	OrderByExprs *OrderByExpression
}

type UnionSelectExpression struct {
	All        *Token
	SelectExpr *SelectExpression
}

type ColumnExpression struct {
	Column *Token
	ArithmeticExpr *ArithmeticExpression
	Function *Token
	Name *Token
}

type JoinExpression struct {
	Type *Token
	Table *Token
	Conditions []*JoinCondition
}

type WhereExpression struct {
	Comparisons []*ComparisonExpression
}

type GroupByExpression struct {
	Columns []*Token
}

type OrderByExpression struct {
	ColumnOrderExprs []*ColumnOrderExpression
}

type JoinCondition struct {
	Left *Token
	Right *Token
}

type ComparisonExpression struct {
	LogicalOperator *Token
	Not *Token
	Left *Token
	Comparator *Token
	Right *Token
}

type ColumnOrderExpression struct {
	Order *Token
	Column *Token
}

type ArithmeticExpression struct {
	Left *Token
	Operator *Token
	Right *Token
}