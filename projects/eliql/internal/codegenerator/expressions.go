package codegenerator

type Expression interface {
	accept(visitor Visitor) Output
}

type Visitor interface {
	visitUnionExpression(u *UnionExpression) Output
	visitSelectExpression(s *SelectExpression) Output
	visitUnionSelectExpression(u *UnionSelectExpression) Output
	visitColumnExpression(c *ColumnExpression) Output
	visitJoinExpression(j *JoinExpression) Output
	visitWhereExpression(w *WhereExpression) Output
	visitGroupByExpression(g *GroupByExpression) Output
	visitOrderByExpression(o *OrderByExpression) Output
	visitJoinCondition(j *JoinCondition) Output
	visitComparisonExpression(c *ComparisonExpression) Output
	visitColumnOrderExpression(c *ColumnOrderExpression) Output
	visitArithmeticExpression(a *ArithmeticExpression) Output
	visitPrimaryExpression(p *PrimaryExpression) Output
}

type Output interface{}

type UnionExpression struct {
	SelectExprs []*UnionSelectExpression
}

func (u *UnionExpression) accept(visitor Visitor) Output {
	return visitor.visitUnionExpression(u)
}

type SelectExpression struct {
	ColumnExprs []*ColumnExpression
	Table       *Token
	JoinExprs   []*JoinExpression
	WhereExpr   *WhereExpression
	GroupByExpr *GroupByExpression
	OrderByExpr *OrderByExpression
}

func (s *SelectExpression) accept(visitor Visitor) Output {
	return visitor.visitSelectExpression(s)
}

type UnionSelectExpression struct {
	All        *Token
	SelectExpr *SelectExpression
}

func (u *UnionSelectExpression) accept(visitor Visitor) Output {
	return visitor.visitUnionSelectExpression(u)
}

type ColumnExpression struct {
	Column         *Token
	ArithmeticExpr *ArithmeticExpression
	Function       *Token
	Name           *Token
}

func (c *ColumnExpression) accept(visitor Visitor) Output {
	return visitor.visitColumnExpression(c)
}

type JoinExpression struct {
	Type       *Token
	Table      *Token
	Conditions []*JoinCondition
}

func (j *JoinExpression) accept(visitor Visitor) Output {
	return visitor.visitJoinExpression(j)
}

type WhereExpression struct {
	Comparisons []*ComparisonExpression
}

func (w *WhereExpression) accept(visitor Visitor) Output {
	return visitor.visitWhereExpression(w)
}

type GroupByExpression struct {
	Columns []*Token
}

func (g *GroupByExpression) accept(visitor Visitor) Output {
	return visitor.visitGroupByExpression(g)
}

type OrderByExpression struct {
	ColumnOrderExprs []*ColumnOrderExpression
}

func (o *OrderByExpression) accept(visitor Visitor) Output {
	return visitor.visitOrderByExpression(o)
}

type JoinCondition struct {
	Left  *Token
	Right *Token
}

func (j *JoinCondition) accept(visitor Visitor) Output {
	return visitor.visitJoinCondition(j)
}

type ComparisonExpression struct {
	LogicalOperator *Token
	Not             *Token
	Left            *Token
	Comparator      *Token
	Right           Expression
}

func (c *ComparisonExpression) accept(visitor Visitor) Output {
	return visitor.visitComparisonExpression(c)
}

type ColumnOrderExpression struct {
	Order  *Token
	Column *Token
}

func (c *ColumnOrderExpression) accept(visitor Visitor) Output {
	return visitor.visitColumnOrderExpression(c)
}

type ArithmeticExpression struct {
	Left     Expression
	Operator *Token
	Right    *PrimaryExpression
}

func (a *ArithmeticExpression) accept(visitor Visitor) Output {
	return visitor.visitArithmeticExpression(a)
}

type PrimaryExpression struct {
	Token *Token
}

func (p *PrimaryExpression) accept(visitor Visitor) Output {
	return visitor.visitPrimaryExpression(p)
}
