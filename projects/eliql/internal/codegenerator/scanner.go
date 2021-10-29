package codegenerator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrEof = errors.New("end of file")
var functionMap = map[string]TokenType{
	"min":         MinFunc,
	"max":         MaxFunc,
	"avg":         AvgFunc,
	"range":       RangeFunc,
	"sum":         SumFunc,
	"count":       CountFunc,
	"now":         NowFunc,
	"to_timezone": ToTimezoneFunc,
	"today":       TodayFunc,
	"interval":    IntervalFunc,
	"concat":      ConcatFunc,
}
var keywordMap = map[string]TokenType{
	"select": Select,
	"from":   From,
	"as":     As,
	"inner":  Inner,
	"left":   Left,
	"right":  Right,
	"full":   Full,
	"join":   Join,
	"on":     On,
	"group":  Group,
	"by":     By,
	"order":  Order,
	"desc":   Desc,
	"asc":    Asc,
	"all":    All,
	"union":  Union,
	"where":  Where,
	"or":     Or,
	"and":    And,
	"not":    Not,
}

type Scanner struct {
	eliql   *Eliql
	source  []rune
	tokens  []*Token
	start   int64
	current int64
	line    int64
}

func NewScanner(eliql *Eliql, source string) *Scanner {
	return &Scanner{eliql: eliql, source: []rune(source), tokens: []*Token{}, line: 1, start: 0, current: 0}
}

func (s *Scanner) ScanTokens() ([]*Token, error) {
	for {
		if s.isAtEnd() {
			break
		}

		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(Eof, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= int64(len(s.source))
}

func (s *Scanner) scanToken() error {
	c, err := s.advance(1)
	if err == ErrEof {
		return nil
	}

	switch c {
	case '(':
		s.addToken(LeftParen, nil)
	case ')':
		s.addToken(RightParen, nil)
	case ',':
		s.addToken(Comma, nil)
	case '+':
		s.addToken(Plus, nil)
	case '/':
		s.addToken(Slash, nil)
	case '*':
		s.addToken(Star, nil)
	case '=':
		s.addToken(Equal, nil)
	case '>':
		if s.matchNext('=') {
			s.addToken(GreaterEqual, nil)
		} else {
			s.addToken(Greater, nil)
		}
	case '<':
		if s.matchNext('=') {
			s.addToken(LessEqual, nil)
		} else {
			s.addToken(Less, nil)
		}
	case '-':
		if s.matchNext('-') {
			for {
				if nextRune, err := s.peek(1); nextRune == '\n' || err == ErrEof {
					break
				}

				s.advance(1)
			}
		} else {
			s.addToken(Minus, nil)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.extractColumnNameOrTableName()
	case '\'':
		s.extractString()
	case ';':
		s.addToken(SemiColon, nil)
	default:
		if s.isDigit(c) {
			s.extractNumber()
		} else if s.isAlpha(c) {
			s.extractKeywordOrFunction()
		} else {
			s.eliql.Error(s.line, "Unexpected character.")
		}
	}
	return nil
}

func (s *Scanner) extractKeywordOrFunction() {
	for isDone := false; !isDone; {
		nextRune, err := s.peek(1)
		if err == ErrEof || s.isAtEnd() {
			isDone = true
			break
		}

		switch nextRune {
		case '(':
			s.extractFunction()
			return
		case ';':
			isDone = true
		case ' ':
			isDone = true
		case '\n':
			isDone = true
		case '\t':
			isDone = true
		default:
			if s.isAlphanumeric(nextRune) {
				_, err := s.advance(1)
				if err == ErrEof {
					isDone = true
				}
			} else {
				s.eliql.Error(s.line, fmt.Sprintf("Unexpected character %v", nextRune))
				return
			}
		}
	}

	word := strings.ToLower(string(s.source[s.start : s.current]))
	keyword, ok := keywordMap[word]
	if !ok {
		s.eliql.Error(s.line, fmt.Sprintf("Unknown keyword '%s'", word))
		return
	}

	s.addToken(keyword, nil)
}

func (s *Scanner) extractFunction() {
	startOfFunc := s.start
	word := strings.ToLower(string(s.source[s.start:s.current]))
	function, ok := functionMap[word]
	if !ok {
		s.eliql.Error(s.line, fmt.Sprintf("Unknown function name '%s'", word))
		return
	}

	s.skipToNext('(', false)

	for {
		c, err := s.advance(1)
		if err == ErrEof {
			s.eliql.Error(s.line, "Unterminated function")
			return
		}

		if c == ')' {
			nestedScanner := NewScanner(s.eliql, string(s.source[s.start:s.current-1]))
			nestedScanner.ScanTokens()

			numberOfTokens := len(nestedScanner.tokens)
			parameters := nestedScanner.tokens[:numberOfTokens-1]
			s.start = startOfFunc
			s.addToken(function, FunctionLiteral{Type: function, Parameters: parameters})
			return
		}
	}
}

func (s *Scanner) extractColumnNameOrTableName() {
	columnLiteral := ColumnLiteral{}
	isColumn := false
	tableStart := s.current
	columnStart := s.current

	for {
		nextRune, err := s.peek(1)
		if err == ErrEof {
			s.eliql.Error(s.line, "Unterminated column name")
			return
		}

		if nextRune == '\n' {
			s.line++
		}

		if nextRune == '"' {
			runeTwoStepsInfront, _ := s.peek(2)
			runeThreeStepsInfront, _ := s.peek(3)

			if runeTwoStepsInfront == '.' && runeThreeStepsInfront == '"' {
				isColumn = true
				columnLiteral.Table = string(s.source[tableStart:s.current])
				s.advance(2)
				columnStart = s.current + 1
			} else if isColumn {
				columnLiteral.Column = string(s.source[columnStart:s.current])
				s.advance(1)
				s.addToken(Column, columnLiteral)
				break
			} else {
				tableName := StringLiteral(s.source[tableStart:s.current])
				s.advance(1)
				s.addToken(Table, tableName)
				break
			}
		}

		s.advance(1)
	}
}

func (s *Scanner) extractString() {
	startIndex := s.current
	endIndex := s.current

	for {
		nextRune, err := s.peek(1)
		if err == ErrEof {
			s.eliql.Error(s.line, "Unterminated string")
			return
		}

		if nextRune == '\n' {
			s.line++
		}

		if nextRune == '\'' {
			endIndex = s.current
			s.advance(1)
			break
		}

		s.advance(1)
	}

	s.addToken(String, StringLiteral(s.source[startIndex:endIndex]))
}

func (s *Scanner) extractNumber() {
	for {
		nextRune, err := s.peek(1)
		if err == ErrEof {
			break
		}

		if s.isDigit(nextRune) {
			s.advance(1)
		} else if runeAfterNext, e := s.peek(2); nextRune == '.' && e == nil && s.isDigit(runeAfterNext) {
			s.advance(2)
		} else {
			break
		}
	}

	num, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		s.eliql.Error(s.line, "A poorly formatted number")
	}

	s.addToken(Number, NumberLiteral(num))
}

func (s *Scanner) peek(step int64) (rune, error) {
	nextIndex := s.current + step - 1

	if s.isAtEnd() || nextIndex >= int64(len(s.source)) {
		return 0, ErrEof
	}

	return s.source[nextIndex], nil
}

func (s *Scanner) advance(step int64) (rune, error) {
	if s.isAtEnd() {
		return 0, ErrEof
	}

	nextRune := s.source[s.current]
	nextIndex := s.current + step
	s.current = nextIndex

	return nextRune, nil
}

func (s *Scanner) addToken(tokenType TokenType, literal TokenLiteral) {
	text := string(s.source[s.start : s.current])
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) matchNext(char rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != char {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isAlphanumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) skipToNext(char rune, inclusive bool) {
	offset := int64(0)
	if inclusive {
		offset = -1
	}

	for {
		c, err := s.advance(1)
		if err == ErrEof {
			s.eliql.Error(s.line, "Unexpected end of file")
			return
		}

		if c == char {
			s.start = s.current + offset
			break
		}
	}
}

func areScannerEqual(expected Scanner, scanner Scanner) bool {
	return expected.eliql == scanner.eliql &&
		string(expected.source) == string(scanner.source) &&
		areTokenSlicesEqual(expected.tokens, scanner.tokens) &&
		expected.start == scanner.start &&
		expected.current == scanner.current &&
		expected.line == scanner.line
}

