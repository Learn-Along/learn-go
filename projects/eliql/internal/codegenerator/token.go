package codegenerator

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal TokenLiteral
	Line    int64
}

// NewToken Creates a New Token
func NewToken(tokenType TokenType, lexeme string, literal TokenLiteral, line int64) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %s %v", t.Type, t.Lexeme, t.Literal)
}

func areTokenSlicesEqual(first []*Token, second []*Token) bool {
	if len(first) != len(second) {
		return false
	}

	for i := range first {
		if !areTokenEqual(*first[i], *second[i]) {
			return false
		}
	}

	return true
}

func areTokenPtrEqual(firstPtr *Token, secondPtr *Token) bool {
	if firstPtr != nil && secondPtr != nil && !areTokenEqual(*firstPtr, *secondPtr) {
		return false
	}

	if (firstPtr == nil || secondPtr == nil) && (firstPtr != secondPtr) {
		return false
	}

	return true
}

func areTokenEqual(first Token, second Token) bool {
	areLiteralsEqual := false
	if first.Literal != nil {
		areLiteralsEqual = first.Literal.Equals(second.Literal)
	} else {
		areLiteralsEqual = first.Literal == second.Literal
	}

	return first.Line == second.Line &&
		first.Type == second.Type &&
		first.Lexeme == second.Lexeme &&
		areLiteralsEqual
}
