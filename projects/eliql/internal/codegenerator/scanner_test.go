package codegenerator

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var testTokens = generateTestTokens()

func TestNewScanner(t *testing.T) {
	ql := &Eliql{}
	source := `SELECT "foo" FROM "bar";`
	sc := NewScanner(ql, source)
	expected := Scanner{
		eliql:   ql,
		source:  []rune(source),
		tokens:  []*Token{},
		start:   0,
		current: 0,
		line:    1,
	}
	assert.True(t, areScannerEqual(expected, *sc))

}

func TestScanner_ScanTokens(t *testing.T) {
	testData := map[string][]*Token{
		`SELECT "bar"."foo" 
         FROM "bar";`: {
			{
				Type:    Select,
				Lexeme:  "SELECT",
				Literal: nil,
				Line:    1,
			},
			{
				Type:   Column,
				Lexeme: `"bar"."foo"`,
				Literal: ColumnLiteral{
					Table:  "bar",
					Column: "foo",
				},
				Line: 1,
			},
			{
				Type:    From,
				Lexeme:  "FROM",
				Literal: nil,
				Line:    2,
			},
			{
				Type:    Table,
				Lexeme:  `"bar"`,
				Literal: StringLiteral("bar"),
				Line:    2,
			},
			{
				Type:    SemiColon,
				Lexeme:  ";",
				Literal: nil,
				Line:    2,
			},
			{
				Type:    Eof,
				Lexeme:  "",
				Literal: nil,
				Line:    2,
			},
		},
	}
	ql := &Eliql{}

	for source, expectedTokens := range testData {
		sc := NewScanner(ql, source)
		sc.ScanTokens()

		assert.True(t,
			areTokenSlicesEqual(expectedTokens, sc.tokens),
			"expected \n%#v or \n%v; got \n%#v or \n%v",
			expectedTokens, expectedTokens, sc.tokens, sc.tokens)
	}
}

func TestScanner_scanToken(t *testing.T) {

	t.Run("Keywords can be extracted", func(t *testing.T) {
		keywordExpectedTokenSliceMap := map[string][]*Token{
			// Select
			`SELECT`:  {testTokens[Select]},
			`SELECT `: {testTokens[Select]},
			`SELECT;`: {testTokens[Select]},

			// From
			`FROM`:  {testTokens[From]},
			`FROM `: {testTokens[From]},
			`FROM;`: {testTokens[From]},

			// As
			`AS`:  {testTokens[As]},
			`AS `: {testTokens[As]},
			`AS;`: {testTokens[As]},

			// Inner
			`INNER`:  {testTokens[Inner]},
			`INNER `: {testTokens[Inner]},
			`INNER;`: {testTokens[Inner]},

			// Left
			`LEFT`:  {testTokens[Left]},
			`LEFT `: {testTokens[Left]},
			`LEFT;`: {testTokens[Left]},

			// Right
			`RIGHT`:  {testTokens[Right]},
			`RIGHT `: {testTokens[Right]},
			`RIGHT;`: {testTokens[Right]},

			// Full
			`FULL`:  {testTokens[Full]},
			`FULL `: {testTokens[Full]},
			`FULL;`: {testTokens[Full]},

			// Join
			`JOIN`:  {testTokens[Join]},
			`JOIN `: {testTokens[Join]},
			`JOIN;`: {testTokens[Join]},

			// On
			`ON`:  {testTokens[On]},
			`ON `: {testTokens[On]},
			`ON;`: {testTokens[On]},

			// Group
			`GROUP`:  {testTokens[Group]},
			`GROUP `: {testTokens[Group]},
			`GROUP;`: {testTokens[Group]},

			// By
			`BY`:  {testTokens[By]},
			`BY `: {testTokens[By]},
			`BY;`: {testTokens[By]},

			// Order
			`ORDER`:  {testTokens[Order]},
			`ORDER `: {testTokens[Order]},
			`ORDER;`: {testTokens[Order]},

			// Desc
			`DESC`:  {testTokens[Desc]},
			`DESC `: {testTokens[Desc]},
			`DESC;`: {testTokens[Desc]},

			// Asc
			`ASC`:  {testTokens[Asc]},
			`ASC `: {testTokens[Asc]},
			`ASC;`: {testTokens[Asc]},

			// All
			`ALL`:  {testTokens[All]},
			`ALL `: {testTokens[All]},
			`ALL;`: {testTokens[All]},

			// Union
			`UNION`:  {testTokens[Union]},
			`UNION `: {testTokens[Union]},
			`UNION;`: {testTokens[Union]},

			//Where
			`WHERE`:  {testTokens[Where]},
			`WHERE `: {testTokens[Where]},
			`WHERE;`: {testTokens[Where]},

			//Or
			`OR`:  {testTokens[Or]},
			`OR `: {testTokens[Or]},
			`OR;`: {testTokens[Or]},

			// And
			`AND`:  {testTokens[And]},
			`AND `: {testTokens[And]},
			`AND;`: {testTokens[And]},

			// Not
			`NOT`:  {testTokens[Not]},
			`NOT `: {testTokens[Not]},
			`NOT;`: {testTokens[Not]},
		}

		for source, expectedTokens := range keywordExpectedTokenSliceMap {
			ql := &Eliql{}
			sc := NewScanner(ql, source)
			sc.scanToken()

			assert.True(t,
				areTokenSlicesEqual(expectedTokens, sc.tokens),
				"expected \n%#v or \n%v; got \n%#v or \n%v",
				expectedTokens, expectedTokens, sc.tokens, sc.tokens)
		}
	})

	t.Run("Functions can be extracted", func(t *testing.T) {
		functionExpectedTokenSliceMap := map[string][]*Token{
			// MIN
			`MIN("foo"."bar")`:  {testTokens[MinFunc]},
			`MIN("foo"."bar") `: {testTokens[MinFunc]},
			`MIN("foo"."bar");`: {testTokens[MinFunc]},

			// MAX
			`MAX("foo"."bar")`:  {testTokens[MaxFunc]},
			`MAX("foo"."bar") `: {testTokens[MaxFunc]},
			`MAX("foo"."bar");`: {testTokens[MaxFunc]},

			// AVG
			`AVG("foo"."bar")`:  {testTokens[AvgFunc]},
			`AVG("foo"."bar") `: {testTokens[AvgFunc]},
			`AVG("foo"."bar");`: {testTokens[AvgFunc]},

			// RANGE
			`RANGE("foo"."bar")`:  {testTokens[RangeFunc]},
			`RANGE("foo"."bar") `: {testTokens[RangeFunc]},
			`RANGE("foo"."bar");`: {testTokens[RangeFunc]},

			// SUM
			`SUM("foo"."bar")`:  {testTokens[SumFunc]},
			`SUM("foo"."bar") `: {testTokens[SumFunc]},
			`SUM("foo"."bar");`: {testTokens[SumFunc]},

			// COUNT
			`COUNT("foo"."bar")`:  {testTokens[CountFunc]},
			`COUNT("foo"."bar") `: {testTokens[CountFunc]},
			`COUNT("foo"."bar");`: {testTokens[CountFunc]},

			// NOW
			`NOW()`:  {testTokens[NowFunc]},
			`NOW() `: {testTokens[NowFunc]},
			`NOW();`: {testTokens[NowFunc]},

			// TO_TIMEZONE
			`TO_TIMEZONE('Africa/Kampala')`:  {testTokens[ToTimezoneFunc]},
			`TO_TIMEZONE('Africa/Kampala') `: {testTokens[ToTimezoneFunc]},
			`TO_TIMEZONE('Africa/Kampala');`: {testTokens[ToTimezoneFunc]},

			// TODAY
			`TODAY()`:  {testTokens[TodayFunc]},
			`TODAY() `: {testTokens[TodayFunc]},
			`TODAY();`: {testTokens[TodayFunc]},

			// CONCAT
			`CONCAT("foo"."bar", '-', "foo"."doe")`:  {testTokens[ConcatFunc]},
			`CONCAT("foo"."bar", '-', "foo"."doe") `: {testTokens[ConcatFunc]},
			`CONCAT("foo"."bar", '-', "foo"."doe");`: {testTokens[ConcatFunc]},

			// INTERVAL
			`INTERVAL('1 day')`:  {testTokens[IntervalFunc]},
			`INTERVAL('1 day') `: {testTokens[IntervalFunc]},
			`INTERVAL('1 day');`: {testTokens[IntervalFunc]},
		}

		for source, expectedTokens := range functionExpectedTokenSliceMap {
			ql := &Eliql{}
			sc := NewScanner(ql, source)
			sc.scanToken()

			assert.True(t,
				areTokenSlicesEqual(expectedTokens, sc.tokens),
				"expected \n%#v or \n%v; got \n%#v or \n%v",
				expectedTokens, expectedTokens, sc.tokens, sc.tokens)
		}
	})

	t.Run("Literals can be extracted", func(t *testing.T) {
		floatNumber := Token{
			Type:    Number,
			Lexeme:  "67.89",
			Literal: NumberLiteral(67.89),
			Line:    1,
		}

		literalExpectedTokenSliceMap := map[string][]*Token{
			// String
			`'foo'`:  {testTokens[String]},
			`'foo' `: {testTokens[String]},
			`'foo';`: {testTokens[String]},

			// Number
			`67`:  {testTokens[Number]},
			`67 `: {testTokens[Number]},
			`67;`: {testTokens[Number]},

			// Float Number
			`67.89`:  {&floatNumber},
			`67.89 `: {&floatNumber},
			`67.89;`: {&floatNumber},

			// Table
			`"foo"`:  {testTokens[Table]},
			`"foo" `: {testTokens[Table]},
			`"foo";`: {testTokens[Table]},

			// Column
			`"foo"."bar"`:  {testTokens[Column]},
			`"foo"."bar" `: {testTokens[Column]},
			`"foo"."bar";`: {testTokens[Column]},
		}

		for source, expectedTokens := range literalExpectedTokenSliceMap {
			ql := &Eliql{}
			sc := NewScanner(ql, source)
			sc.scanToken()

			assert.True(t,
				areTokenSlicesEqual(expectedTokens, sc.tokens),
				"expected \n%#v or \n%v; got \n%#v or \n%v",
				expectedTokens, expectedTokens, sc.tokens, sc.tokens)
		}
	})

	t.Run("Comparators can be extracted", func(t *testing.T) {
		comparatorExpectedTokenSliceMap := map[string][]*Token{
			// Equal
			`=`:  {testTokens[Equal]},
			`= `: {testTokens[Equal]},
			`=;`: {testTokens[Equal]},

			// Greater
			`>`:  {testTokens[Greater]},
			`> `: {testTokens[Greater]},
			`>;`: {testTokens[Greater]},

			// GreaterEqual
			`>=`:  {testTokens[GreaterEqual]},
			`>= `: {testTokens[GreaterEqual]},
			`>=;`: {testTokens[GreaterEqual]},

			// Less
			`<`:  {testTokens[Less]},
			`< `: {testTokens[Less]},
			`<;`: {testTokens[Less]},

			// LessEqual
			`<=`:  {testTokens[LessEqual]},
			`<= `: {testTokens[LessEqual]},
			`<=;`: {testTokens[LessEqual]},
		}

		for source, expectedTokens := range comparatorExpectedTokenSliceMap {
			ql := &Eliql{}
			sc := NewScanner(ql, source)
			sc.scanToken()

			assert.True(t,
				areTokenSlicesEqual(expectedTokens, sc.tokens),
				"expected \n%#v or \n%v; got \n%#v or \n%v",
				expectedTokens, expectedTokens, sc.tokens, sc.tokens)
		}
	})

	t.Run("Symbols can be extracted", func(t *testing.T) {
		symbolsExpectedTokenSliceMap := map[string][]*Token{
			// LeftParen
			`(`:  {testTokens[LeftParen]},
			`( `: {testTokens[LeftParen]},
			`(;`: {testTokens[LeftParen]},

			// RightParen
			`)`:  {testTokens[RightParen]},
			`) `: {testTokens[RightParen]},
			`);`: {testTokens[RightParen]},

			// Comma
			`,`:  {testTokens[Comma]},
			`, `: {testTokens[Comma]},
			`,;`: {testTokens[Comma]},

			// Minus
			`-`:  {testTokens[Minus]},
			`- `: {testTokens[Minus]},
			`-;`: {testTokens[Minus]},

			// Plus
			`+`:  {testTokens[Plus]},
			`+ `: {testTokens[Plus]},
			`+;`: {testTokens[Plus]},

			// Slash
			`/`:  {testTokens[Slash]},
			`/ `: {testTokens[Slash]},
			`/;`: {testTokens[Slash]},

			// Star
			`*`:  {testTokens[Star]},
			`* `: {testTokens[Star]},
			`*;`: {testTokens[Star]},

			// SemiColon
			`;`:  {testTokens[SemiColon]},
			`; `: {testTokens[SemiColon]},
			`;;`: {testTokens[SemiColon]},
		}

		for source, expectedTokens := range symbolsExpectedTokenSliceMap {
			ql := &Eliql{}
			sc := NewScanner(ql, source)
			sc.scanToken()

			assert.True(t,
				areTokenSlicesEqual(expectedTokens, sc.tokens),
				"expected \n%#v or \n%v; got \n%#v or \n%v",
				expectedTokens, expectedTokens, sc.tokens, sc.tokens)
		}
	})
}

func TestScanner_extractFunction(t *testing.T) {
	type testData struct {
		err             error
		tokens          []*Token
		expectedCurrent int64
	}

	sourceTestDataMap := map[string]testData{
		// MIN
		`MIN("foo"."bar")`: {
			err:             nil,
			tokens:          []*Token{testTokens[MinFunc]},
			expectedCurrent: 16,
		},
		`MIN("foo"."bar") `: {
			err:             nil,
			tokens:          []*Token{testTokens[MinFunc]},
			expectedCurrent: 16,
		},
		`MIN("foo"."bar");`: {
			err:             nil,
			tokens:          []*Token{testTokens[MinFunc]},
			expectedCurrent: 16,
		},

		// MAX
		`MAX("foo"."bar")`: {
			err:             nil,
			tokens:          []*Token{testTokens[MaxFunc]},
			expectedCurrent: 16,
		},
		`MAX("foo"."bar") `: {
			err:             nil,
			tokens:          []*Token{testTokens[MaxFunc]},
			expectedCurrent: 16,
		},
		`MAX("foo"."bar");`: {
			err:             nil,
			tokens:          []*Token{testTokens[MaxFunc]},
			expectedCurrent: 16,
		},

		//AVG
		`AVG("foo"."bar")`: {
			err:             nil,
			tokens:          []*Token{testTokens[AvgFunc]},
			expectedCurrent: 16,
		},
		`AVG("foo"."bar") `: {
			err:             nil,
			tokens:          []*Token{testTokens[AvgFunc]},
			expectedCurrent: 16,
		},
		`AVG("foo"."bar");`: {
			err:             nil,
			tokens:          []*Token{testTokens[AvgFunc]},
			expectedCurrent: 16,
		},

		//RANGE
		`RANGE("foo"."bar")`: {
			err:             nil,
			tokens:          []*Token{testTokens[RangeFunc]},
			expectedCurrent: 18,
		},
		`RANGE("foo"."bar") `: {
			err:             nil,
			tokens:          []*Token{testTokens[RangeFunc]},
			expectedCurrent: 18,
		},
		`RANGE("foo"."bar");`: {
			err:             nil,
			tokens:          []*Token{testTokens[RangeFunc]},
			expectedCurrent: 18,
		},

		// SUM
		`SUM("foo"."bar")`: {
			err:             nil,
			tokens:          []*Token{testTokens[SumFunc]},
			expectedCurrent: 16,
		},
		`SUM("foo"."bar") `: {
			err:             nil,
			tokens:          []*Token{testTokens[SumFunc]},
			expectedCurrent: 16,
		},
		`SUM("foo"."bar");`: {
			err:             nil,
			tokens:          []*Token{testTokens[SumFunc]},
			expectedCurrent: 16,
		},

		// COUNT
		`COUNT("foo"."bar")`: {
			err:             nil,
			tokens:          []*Token{testTokens[CountFunc]},
			expectedCurrent: 18,
		},
		`COUNT("foo"."bar") `: {
			err:             nil,
			tokens:          []*Token{testTokens[CountFunc]},
			expectedCurrent: 18,
		},
		`COUNT("foo"."bar");`: {
			err:             nil,
			tokens:          []*Token{testTokens[CountFunc]},
			expectedCurrent: 18,
		},

		// NOW
		`NOW()`: {
			err:             nil,
			tokens:          []*Token{testTokens[NowFunc]},
			expectedCurrent: 5,
		},
		`NOW() `: {
			err:             nil,
			tokens:          []*Token{testTokens[NowFunc]},
			expectedCurrent: 5,
		},
		`NOW();`: {
			err:             nil,
			tokens:          []*Token{testTokens[NowFunc]},
			expectedCurrent: 5,
		},

		// TO_TIMEZONE
		`TO_TIMEZONE('Africa/Kampala')`: {
			err:             nil,
			tokens:          []*Token{testTokens[ToTimezoneFunc]},
			expectedCurrent: 29,
		},
		`TO_TIMEZONE('Africa/Kampala') `: {
			err:             nil,
			tokens:          []*Token{testTokens[ToTimezoneFunc]},
			expectedCurrent: 29,
		},
		`TO_TIMEZONE('Africa/Kampala');`: {
			err:             nil,
			tokens:          []*Token{testTokens[ToTimezoneFunc]},
			expectedCurrent: 29,
		},

		// TODAY
		`TODAY()`: {
			err:             nil,
			tokens:          []*Token{testTokens[TodayFunc]},
			expectedCurrent: 7,
		},
		`TODAY() `: {
			err:             nil,
			tokens:          []*Token{testTokens[TodayFunc]},
			expectedCurrent: 7,
		},
		`TODAY();`: {
			err:             nil,
			tokens:          []*Token{testTokens[TodayFunc]},
			expectedCurrent: 7,
		},

		// INTERVAL
		`INTERVAL('1 day')`: {
			err:             nil,
			tokens:          []*Token{testTokens[IntervalFunc]},
			expectedCurrent: 17,
		},
		`INTERVAL('1 day') `: {
			err:             nil,
			tokens:          []*Token{testTokens[IntervalFunc]},
			expectedCurrent: 17,
		},
		`INTERVAL('1 day');`: {
			err:             nil,
			tokens:          []*Token{testTokens[IntervalFunc]},
			expectedCurrent: 17,
		},

		//ConcatFunc
		`CONCAT("foo"."bar", '-', "foo"."doe")`: {
			err:             nil,
			tokens:          []*Token{testTokens[ConcatFunc]},
			expectedCurrent: 37,
		},
		`CONCAT("foo"."bar", '-', "foo"."doe") `: {
			err:             nil,
			tokens:          []*Token{testTokens[ConcatFunc]},
			expectedCurrent: 37,
		},
		`CONCAT("foo"."bar", '-', "foo"."doe");`: {
			err:             nil,
			tokens:          []*Token{testTokens[ConcatFunc]},
			expectedCurrent: 37,
		},
	}

	for source, testDatum := range sourceTestDataMap {
		ql := &Eliql{}
		sc := NewScanner(ql, source)
		sc.current = int64(strings.Index(source, "("))
		sc.extractFunction()

		assert.Equal(t, testDatum.expectedCurrent, sc.current)
		assert.True(t,
			areTokenSlicesEqual(testDatum.tokens, sc.tokens),
			"expected \n%#v or \n%v; got \n%#v or \n%v",
			testDatum.tokens, testDatum.tokens, sc.tokens, sc.tokens)
	}
}

func TestScanner_extractColumnNameOrTableName(t *testing.T) {
	type testData struct {
		err             error
		tokens          []*Token
		expectedCurrent int64
	}

	sourceTestDataMap := map[string]testData{
		// Column
		// FIXME: Add data for erroneous data, and add custom errors for unclosed table
		`"foo"."bar"`: {
			err:             nil,
			tokens:          []*Token{testTokens[Column]},
			expectedCurrent: 11,
		},
		`"foo"."bar"  `: {
			err:             nil,
			tokens:          []*Token{testTokens[Column]},
			expectedCurrent: 11,
		},
		`"foo"."bar";`: {
			err:             nil,
			tokens:          []*Token{testTokens[Column]},
			expectedCurrent: 11,
		},

		// Table
		`"foo"`: {
			err:             nil,
			tokens:          []*Token{testTokens[Table]},
			expectedCurrent: 5,
		},
		`"foo" `: {
			err:             nil,
			tokens:          []*Token{testTokens[Table]},
			expectedCurrent: 5,
		},
		`"foo";`: {
			err:             nil,
			tokens:          []*Token{testTokens[Table]},
			expectedCurrent: 5,
		},
	}

	for source, testDatum := range sourceTestDataMap {
		ql := &Eliql{}
		sc := NewScanner(ql, source)
		sc.current = int64(strings.Index(source, `"`)) + 1
		sc.extractColumnNameOrTableName()

		assert.Equal(t, testDatum.expectedCurrent, sc.current)
		assert.True(t,
			areTokenSlicesEqual(testDatum.tokens, sc.tokens),
			"expected \n%#v or \n%v; got \n%#v or \n%v",
			testDatum.tokens, testDatum.tokens, sc.tokens, sc.tokens)
	}
}

func TestScanner_extractString(t *testing.T) {
	type testData struct {
		err             error
		tokens          []*Token
		expectedCurrent int64
	}

	sourceTestDataMap := map[string]testData{
		// FIXME: Add data for erroneous data, and add custom errors for unclosed string
		`'foo'`: {
			err:             nil,
			tokens:          []*Token{testTokens[String]},
			expectedCurrent: 5,
		},
		`'foo'  `: {
			err:             nil,
			tokens:          []*Token{testTokens[String]},
			expectedCurrent: 5,
		},
		`'foo';`: {
			err:             nil,
			tokens:          []*Token{testTokens[String]},
			expectedCurrent: 5,
		},
	}

	for source, testDatum := range sourceTestDataMap {
		ql := &Eliql{}
		sc := NewScanner(ql, source)
		sc.current = int64(strings.Index(source, `'`)) + 1
		sc.extractString()

		assert.Equal(t, testDatum.expectedCurrent, sc.current)
		assert.True(t,
			areTokenSlicesEqual(testDatum.tokens, sc.tokens),
			"expected \n%#v or \n%v; got \n%#v or \n%v",
			testDatum.tokens, testDatum.tokens, sc.tokens, sc.tokens)
	}
}

func TestScanner_extractNumber(t *testing.T) {
	type testData struct {
		err             error
		tokens          []*Token
		expectedCurrent int64
	}

	floatNumber := Token{
		Type:    Number,
		Lexeme:  "67.89",
		Literal: NumberLiteral(67.89),
		Line:    1,
	}

	sourceTestDataMap := map[string]testData{
		// FIXME: Add data for erroneous data
		// Integer
		`67`: {
			err:             nil,
			tokens:          []*Token{testTokens[Number]},
			expectedCurrent: 2,
		},
		`67  `: {
			err:             nil,
			tokens:          []*Token{testTokens[Number]},
			expectedCurrent: 2,
		},
		`67;`: {
			err:             nil,
			tokens:          []*Token{testTokens[Number]},
			expectedCurrent: 2,
		},

		// Float
		`67.89`: {
			err:             nil,
			tokens:          []*Token{&floatNumber},
			expectedCurrent: 5,
		},
		`67.89  `: {
			err:             nil,
			tokens:          []*Token{&floatNumber},
			expectedCurrent: 5,
		},
		`67.89;`: {
			err:             nil,
			tokens:          []*Token{&floatNumber},
			expectedCurrent: 5,
		},
	}

	for source, testDatum := range sourceTestDataMap {
		ql := &Eliql{}
		sc := NewScanner(ql, source)
		sc.extractNumber()

		assert.Equal(t, testDatum.expectedCurrent, sc.current)
		assert.True(t,
			areTokenSlicesEqual(testDatum.tokens, sc.tokens),
			"expected \n%#v or \n%v; got \n%#v or \n%v",
			testDatum.tokens, testDatum.tokens, sc.tokens, sc.tokens)
	}
}

func TestScanner_advance(t *testing.T) {
	type testRecord struct {
		step            int64
		err             error
		nextRune        rune
		expectedCurrent int64
	}

	source := "SELECT"
	ql := &Eliql{}
	sc := NewScanner(ql, source)

	testData := []testRecord{
		{step: 1, err: nil, nextRune: 'S', expectedCurrent: 1},
		{step: 3, err: nil, nextRune: 'E', expectedCurrent: 4},
		// no advancing as 2 steps take it out of bounds of the slice
		{step: 2, err: nil, nextRune: 'C', expectedCurrent: 6},
		{step: 0, err: ErrEof, nextRune: 0, expectedCurrent: 6},
	}

	for _, testDatum := range testData {
		actualNextRune, actualError := sc.advance(testDatum.step)

		assert.Equal(t, testDatum.expectedCurrent, sc.current)
		assert.Equal(t, testDatum.nextRune, actualNextRune)
		assert.Equal(t, testDatum.err, actualError)
	}
}

func TestScanner_peek(t *testing.T) {
	type testData struct {
		err             error
		nextRune        rune
		expectedCurrent int64
	}

	source := "SELECT"
	ql := &Eliql{}
	sc := NewScanner(ql, source)

	stepCurrentMap := map[int64]testData{
		1: {err: nil, nextRune: 'S', expectedCurrent: 0},
		2: {err: nil, nextRune: 'E', expectedCurrent: 0},
		3: {err: nil, nextRune: 'L', expectedCurrent: 0},
		4: {err: nil, nextRune: 'E', expectedCurrent: 0},
		5: {err: nil, nextRune: 'C', expectedCurrent: 0},
		6: {err: nil, nextRune: 'T', expectedCurrent: 0},
		//// no advancing as 7 steps take it out of bounds
		7: {err: ErrEof, nextRune: 0, expectedCurrent: 0},
	}

	for step, testDatum := range stepCurrentMap {
		actualNextRune, actualError := sc.peek(step)

		assert.Equal(t, testDatum.expectedCurrent, sc.current)
		assert.Equal(t, testDatum.nextRune, actualNextRune)
		assert.Equal(t, testDatum.err, actualError)
	}
}

func generateTestTokens() map[TokenType]*Token {
	return map[TokenType]*Token{
		// Keywords
		Select: {
			Type:    Select,
			Lexeme:  "SELECT",
			Literal: nil,
			Line:    1,
		},
		From: {
			Type:    From,
			Lexeme:  "FROM",
			Literal: nil,
			Line:    1,
		},
		As: {
			Type:    As,
			Lexeme:  "AS",
			Literal: nil,
			Line:    1,
		},
		Inner: {
			Type:    Inner,
			Lexeme:  "INNER",
			Literal: nil,
			Line:    1,
		},
		Left: {
			Type:    Left,
			Lexeme:  "LEFT",
			Literal: nil,
			Line:    1,
		},
		Right: {
			Type:    Right,
			Lexeme:  "RIGHT",
			Literal: nil,
			Line:    1,
		},
		Full: {
			Type:    Full,
			Lexeme:  "FULL",
			Literal: nil,
			Line:    1,
		},
		Join: {
			Type:    Join,
			Lexeme:  "JOIN",
			Literal: nil,
			Line:    1,
		},
		On: {
			Type:    On,
			Lexeme:  "ON",
			Literal: nil,
			Line:    1,
		},
		Order: {
			Type:    Order,
			Lexeme:  "ORDER",
			Literal: nil,
			Line:    1,
		},
		By: {
			Type:    By,
			Lexeme:  "BY",
			Literal: nil,
			Line:    1,
		},
		Desc: {
			Type:    Desc,
			Lexeme:  "DESC",
			Literal: nil,
			Line:    1,
		},
		Asc: {
			Type:    Asc,
			Lexeme:  "ASC",
			Literal: nil,
			Line:    1,
		},
		All: {
			Type:    All,
			Lexeme:  "ALL",
			Literal: nil,
			Line:    1,
		},
		Union: {
			Type:    Union,
			Lexeme:  "UNION",
			Literal: nil,
			Line:    1,
		},
		Where: {
			Type:    Where,
			Lexeme:  "WHERE",
			Literal: nil,
			Line:    1,
		},
		Or: {
			Type:    Or,
			Lexeme:  "OR",
			Literal: nil,
			Line:    1,
		},
		And: {
			Type:    And,
			Lexeme:  "AND",
			Literal: nil,
			Line:    1,
		},
		Not: {
			Type:    Not,
			Lexeme:  "NOT",
			Literal: nil,
			Line:    1,
		},
		Group: {
			Type:    Group,
			Lexeme:  "GROUP",
			Literal: nil,
			Line:    1,
		},

		// Functions
		MinFunc: {
			Type:   MinFunc,
			Lexeme: `MIN("foo"."bar")`,
			Literal: FunctionLiteral{
				Name: "min",
				Type: MinFunc,
				Parameters: []*Token{
					{
						Type:   Column,
						Lexeme: `"foo"."bar"`,
						Literal: ColumnLiteral{
							Table:  "foo",
							Column: "bar",
						},
						Line: 1,
					},
				},
			},
			Line: 1,
		},
		MaxFunc: {
			Type:   MaxFunc,
			Lexeme: `MAX("foo"."bar")`,
			Literal: FunctionLiteral{
				Name: "max",
				Type: MaxFunc,
				Parameters: []*Token{
					{
						Type:   Column,
						Lexeme: `"foo"."bar"`,
						Literal: ColumnLiteral{
							Table:  "foo",
							Column: "bar",
						},
						Line: 1,
					},
				},
			},
			Line: 1,
		},
		AvgFunc: {
			Type:   AvgFunc,
			Lexeme: `AVG("foo"."bar")`,
			Literal: FunctionLiteral{
				Name: "avg",
				Type: AvgFunc,
				Parameters: []*Token{
					{
						Type:   Column,
						Lexeme: `"foo"."bar"`,
						Literal: ColumnLiteral{
							Table:  "foo",
							Column: "bar",
						},
						Line: 1,
					},
				},
			},
			Line: 1,
		},
		RangeFunc: {
			Type:   RangeFunc,
			Lexeme: `RANGE("foo"."bar")`,
			Literal: FunctionLiteral{
				Name: "range",
				Type: RangeFunc,
				Parameters: []*Token{
					{
						Type:   Column,
						Lexeme: `"foo"."bar"`,
						Literal: ColumnLiteral{
							Table:  "foo",
							Column: "bar",
						},
						Line: 1,
					},
				},
			},
			Line: 1,
		},
		SumFunc: {
			Type:   SumFunc,
			Lexeme: `SUM("foo"."bar")`,
			Literal: FunctionLiteral{
				Name: "sum",
				Type: SumFunc,
				Parameters: []*Token{
					{
						Type:   Column,
						Lexeme: `"foo"."bar"`,
						Literal: ColumnLiteral{
							Table:  "foo",
							Column: "bar",
						},
						Line: 1,
					},
				},
			},
			Line: 1,
		},
		CountFunc: {
			Type:   CountFunc,
			Lexeme: `COUNT("foo"."bar")`,
			Literal: FunctionLiteral{
				Name: "count",
				Type: CountFunc,
				Parameters: []*Token{
					{
						Type:   Column,
						Lexeme: `"foo"."bar"`,
						Literal: ColumnLiteral{
							Table:  "foo",
							Column: "bar",
						},
						Line: 1,
					},
				},
			},
			Line: 1,
		},
		NowFunc: {
			Type:   NowFunc,
			Lexeme: `NOW()`,
			Literal: FunctionLiteral{
				Name:       "now",
				Type:       NowFunc,
				Parameters: []*Token{},
			},
			Line: 1,
		},
		ToTimezoneFunc: {
			Type:   ToTimezoneFunc,
			Lexeme: `TO_TIMEZONE('Africa/Kampala')`,
			Literal: FunctionLiteral{
				Name: "to_timezone",
				Type: ToTimezoneFunc,
				Parameters: []*Token{
					{
						Type:    String,
						Lexeme:  `'Africa/Kampala'`,
						Literal: StringLiteral("Africa/Kampala"),
						Line:    1,
					},
				},
			},
			Line: 1,
		},
		TodayFunc: {
			Type:   TodayFunc,
			Lexeme: `TODAY()`,
			Literal: FunctionLiteral{
				Name:       "today",
				Type:       TodayFunc,
				Parameters: []*Token{},
			},
			Line: 1,
		},
		ConcatFunc: {
			Type:   ConcatFunc,
			Lexeme: `CONCAT("foo"."bar", '-', "foo"."doe")`,
			Literal: FunctionLiteral{
				Name: "concat",
				Type: ConcatFunc,
				Parameters: []*Token{
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
						Type:    Comma,
						Lexeme:  `,`,
						Literal: nil,
						Line:    1,
					},
					{
						Type:    String,
						Lexeme:  `'-'`,
						Literal: StringLiteral("-"),
						Line:    1,
					},
					{
						Type:    Comma,
						Lexeme:  `,`,
						Literal: nil,
						Line:    1,
					},
					{
						Type:   Column,
						Lexeme: `"foo"."doe"`,
						Literal: ColumnLiteral{
							Table:  "foo",
							Column: "doe",
						},
						Line: 1,
					},
				},
			},
			Line: 1,
		},
		IntervalFunc: {
			Type:   IntervalFunc,
			Lexeme: `INTERVAL('1 day')`,
			Literal: FunctionLiteral{
				Name: "interval",
				Type: IntervalFunc,
				Parameters: []*Token{
					{
						Type:    String,
						Lexeme:  `'1 day'`,
						Literal: StringLiteral("1 day"),
						Line:    1,
					},
				},
			},
			Line: 1,
		},

		// Literals
		String: {
			Type:    String,
			Lexeme:  "'foo'",
			Literal: StringLiteral("foo"),
			Line:    1,
		},
		Number: {
			Type:    Number,
			Lexeme:  "67",
			Literal: NumberLiteral(67),
			Line:    1,
		},
		Table: {
			Type:    Table,
			Lexeme:  `"foo"`,
			Literal: StringLiteral("foo"),
			Line:    1,
		},
		Column: {
			Type:   Column,
			Lexeme: `"foo"."bar"`,
			Literal: ColumnLiteral{
				Table:  "foo",
				Column: "bar",
			},
			Line: 1,
		},

		// Comparators
		Equal: {
			Type:    Equal,
			Lexeme:  "=",
			Literal: nil,
			Line:    1,
		},
		Greater: {
			Type:    Greater,
			Lexeme:  ">",
			Literal: nil,
			Line:    1,
		},
		GreaterEqual: {
			Type:    GreaterEqual,
			Lexeme:  ">=",
			Literal: nil,
			Line:    1,
		},
		Less: {
			Type:    Less,
			Lexeme:  "<",
			Literal: nil,
			Line:    1,
		},
		LessEqual: {
			Type:    LessEqual,
			Lexeme:  "<=",
			Literal: nil,
			Line:    1,
		},

		// Symbols
		LeftParen: {
			Type:    LeftParen,
			Lexeme:  "(",
			Literal: nil,
			Line:    1,
		},
		RightParen: {
			Type:    RightParen,
			Lexeme:  ")",
			Literal: nil,
			Line:    1,
		},
		Comma: {
			Type:    Comma,
			Lexeme:  ",",
			Literal: nil,
			Line:    1,
		},
		Minus: {
			Type:    Minus,
			Lexeme:  "-",
			Literal: nil,
			Line:    1,
		},
		Plus: {
			Type:    Plus,
			Lexeme:  "+",
			Literal: nil,
			Line:    1,
		},
		Slash: {
			Type:    Slash,
			Lexeme:  "/",
			Literal: nil,
			Line:    1,
		},
		Star: {
			Type:    Star,
			Lexeme:  "*",
			Literal: nil,
			Line:    1,
		},
		SemiColon: {
			Type:    SemiColon,
			Lexeme:  ";",
			Literal: nil,
			Line:    1,
		},
	}
}
