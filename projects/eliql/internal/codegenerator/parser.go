package codegenerator

/*
Parser converts sequence of tokens into EliQL expressions
	The grammar is as follows:

	expression              -> unionExpr | selectExpr;
	unionExpr               -> selectExpr ("UNION" "ALL"? selectExpr)+;
	selectExpr              -> "SELECT" columnExpr ("," columnExpr)* "FROM" TABLE joinExpr* whereExpr? groupByExpr? orderByExpr?;
	columnExpr              -> (arithmeticExpr | COLUMN | FUNCTION) "AS" NAME;
	joinExpr                -> ("LEFT" | "RIGHT" | "FULL")? "JOIN" TABLE "ON" joinConditionExpr ("AND" joinConditionExpr)*;
	whereExpr               -> "WHERE" comparisonExpr (("AND" | "OR") comparisonExpr)*;
	groupByExpr             -> "GROUP" "BY" COLUMN ("," COLUMN)*;
	orderByExpr             -> "ORDER" "BY" columnOrderExpr ("," columnOrderExpr)*;
	arithmeticExpr          -> "(" primaryExpr ( ("/" | "-" | "+" | "*") primaryExpr )+ ")";
	joinConditionExpr       -> COLUMN "=" COLUMN
	comparisonExpr          -> "NOT"? Column (">" | ">=" | "=" | "<" | "<=") (primaryExpr | arithmeticExpr);
	columnOrderExpr         -> COLUMN ("ASC" | "DESC")
	primaryExpr             -> COLUMN | NUMBER | STRING

*/
type Parser struct {
	tokens  []Token
	current int
}

// NewParser creates a new Parser with given Tokens
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) expression() Expression {
	expr := p.selectExpr()

	if p.check(Union) {
		return p.unionExpr(expr)
	}

	return expr
}

func (p *Parser) unionExpr(left Expression) *UnionExpression {
	expr := &UnionExpression{
		SelectExprs: []*UnionSelectExpression{{
			SelectExpr: left.(*SelectExpression),
		}},
	}

	for p.match(Union) {
		unionSelect := &UnionSelectExpression{}

		if p.match(All) {
			prev := p.previous()
			unionSelect.All = &prev
		}

		unionSelect.SelectExpr = p.selectExpr()
		expr.SelectExprs = append(expr.SelectExprs, unionSelect)
	}

	return expr
}

func (p *Parser) selectExpr() *SelectExpression {
	expr := &SelectExpression{
		ColumnExprs: []*ColumnExpression{},
		Table:       nil,
		JoinExprs:   []*JoinExpression{},
		WhereExpr:   nil,
		GroupByExpr: nil,
		OrderByExpr: nil,
	}

	if p.match(Select) {
		//for colExpr := p.columnExpr(); p.match(Comma); colExpr = p.columnExpr() {
		//	expr.ColumnExprs = append(expr.ColumnExprs, colExpr)
		//}

		for {
			colExpr := p.columnExpr()
			expr.ColumnExprs = append(expr.ColumnExprs, colExpr)

			if !p.match(Comma) {
				break
			}
		}

		if p.match(From) {
			table := p.advance()
			expr.Table = &table
		}

		for j := p.joinExpr(); p.match(Join); j = p.joinExpr() {
			expr.JoinExprs = append(expr.JoinExprs, j)
		}

		if p.match(Where) {
			expr.WhereExpr = p.whereExpr()
		}

		if p.match(Group) {
			expr.GroupByExpr = p.groupByExpr()
		}

		if p.match(Order) {
			expr.OrderByExpr = p.orderByExpr()
		}
	}

	return expr
}

func (p *Parser) columnExpr() *ColumnExpression {
	functionTokenTypes := []TokenType{
		MinFunc,
		MaxFunc,
		AvgFunc,
		RangeFunc,
		SumFunc,
		CountFunc,
		NowFunc,
		ToTimezoneFunc,
		TodayFunc,
		IntervalFunc,
		ConcatFunc,
	}
	expr := &ColumnExpression{}

	if p.match(functionTokenTypes...) {
		function := p.previous()
		expr.Function = &function
	}

	if p.match(Column) {
		col := p.previous()
		expr.Column = &col
	}

	if p.match(LeftParen) {
		expr.ArithmeticExpr = p.arithmeticExpr()
	}

	if p.match(As) {
		name := p.advance()
		expr.Name = &name
	}

	return expr
}

func (p *Parser) joinExpr() *JoinExpression {
	permittedTokenTypes := []TokenType{
		Left,
		Right,
		Full,
	}
	expr := &JoinExpression{
		Conditions: []*JoinCondition{},
	}

	tokenBeforeJoin := p.backPeek(2)
	for _, tokenType := range permittedTokenTypes {
		if tokenBeforeJoin.Type == tokenType {
			expr.Type = &tokenBeforeJoin
		}
	}

	if p.match(Table) {
		table := p.previous()
		expr.Table = &table
	}

	if p.match(On) {
		for p.match(And) {
			expr.Conditions = append(expr.Conditions, p.joinConditionExpr())
		}
	}

	return expr
}

func (p *Parser) whereExpr() *WhereExpression {
	expr := &WhereExpression{
		Comparisons: []*ComparisonExpression{p.comparisonExpr(nil)},
	}

	for p.match(And, Or) {
		logicalOperator := p.previous()
		comp := p.comparisonExpr(&logicalOperator)
		expr.Comparisons = append(expr.Comparisons, comp)
	}

	return expr
}

func (p *Parser) groupByExpr() *GroupByExpression {
	expr := &GroupByExpression{
		Columns: []*Token{},
	}

	if p.match(By) {
		for col := p.advance(); p.match(Comma); col = p.advance() {
			expr.Columns = append(expr.Columns, &col)
		}
	}

	return expr
}

func (p *Parser) orderByExpr() *OrderByExpression {
	expr := &OrderByExpression{
		ColumnOrderExprs: []*ColumnOrderExpression{},
	}

	if p.match(By) {
		for colOrder := p.columnOrderExpr(); p.match(Comma); colOrder = p.columnOrderExpr() {
			expr.ColumnOrderExprs = append(expr.ColumnOrderExprs, colOrder)
		}
	}

	return expr
}

func (p *Parser) arithmeticExpr() *ArithmeticExpression {
	expr := &ArithmeticExpression{}
	expr.Left = p.primaryExpr()

	for p.match(Slash, Star, Plus, Minus) {
		operator := p.previous()
		right := p.primaryExpr()

		if expr.Right != nil {
			expr.Left = expr
		}

		expr.Operator = &operator
		expr.Right = right
	}

	return expr
}

func (p *Parser) joinConditionExpr() *JoinCondition {
	expr := &JoinCondition{}

	if p.match(Column) {
		col := p.previous()
		expr.Left = &col
	}

	if p.match(Equal) {
		col := p.advance()
		expr.Right = &col
	}

	return expr
}

func (p *Parser) comparisonExpr(logicalOperator *Token) *ComparisonExpression {
	expr := &ComparisonExpression{
		LogicalOperator: logicalOperator,
	}

	if p.match(Not) {
		not := p.previous()
		expr.Not = &not
	}

	if p.match(Column) {
		col := p.previous()
		expr.Left = &col
	}

	if p.match(Greater, Less, Equal, GreaterEqual, LessEqual) {
		comp := p.previous()
		expr.Comparator = &comp
	}

	if p.match(LeftParen) {
		expr.Right = p.arithmeticExpr()
	} else {
		expr.Right = p.primaryExpr()
	}

	return expr
}

func (p *Parser) columnOrderExpr() *ColumnOrderExpression {
	expr := &ColumnOrderExpression{}

	if p.match(Column) {
		col := p.previous()
		expr.Column = &col
	}

	if p.match(Order) {
		order := p.previous()
		expr.Order = &order
	}

	return expr
}

func (p *Parser) primaryExpr() *PrimaryExpression {
	expr := &PrimaryExpression{}

	if p.match(Column, Number, String) {
		tkn := p.previous()
		expr.Token = &tkn
	}

	return expr
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current += 1
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == Eof
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) backPeek(places int) Token {
	return p.tokens[p.current-places]
}
