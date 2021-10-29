package codegenerator

type TokenType int64

const (
	LeftParen TokenType = iota
	RightParen
	Comma
	Minus
	Plus
	Slash
	Star
	SemiColon

	Equal
	Greater
	GreaterEqual
	Less
	LessEqual

	String
	Number
	Table
	Column

	Select
	From
	As
	Inner
	Left
	Right
	Full
	Join
	On
	Group
	By
	Order
	Desc
	Asc
	All
	Union
	Where
	Or
	And
	Not

	MinFunc
	MaxFunc
	AvgFunc
	RangeFunc
	SumFunc
	CountFunc
	NowFunc
	ToTimezoneFunc
	TodayFunc
	IntervalFunc
	ConcatFunc

	Eof
)
