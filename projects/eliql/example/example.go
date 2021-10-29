package main

import (
	//"fmt"
	//"regexp"
	"github.com/Learn-Along/learn-go/projects/eliql/internal/codegenerator"
)

func main() {
	//word := `[0-9a-zA-Z/\\ \-\*\+\.=_]+`
	//stringLiteral := fmt.Sprintf(`'?%s'?`, word)
	//timeInstanceFunc := `(?:NOW\(\)|TODAY\(\))`
	//durationLiteral := `(?:(\d+\s+years?)?(\s+\d+\s+days?)?)(\s+\d+\s+hours?)?(\s+\d+\s+minutes?)?(\s+\d+\s+seconds?)?(\s+\d+\s+milliseconds?)?`
	//intervalFuncLiteral := fmt.Sprintf(`(?:%s\s+(?:\-|\+)\s+INTERVAL\(\s*'%s'\s*\))`, timeInstanceFunc, durationLiteral)
	//timeLiteral := fmt.Sprintf(`(?:%s|%s|TO_TIMEZONE\(\s*(?:%s|%s),\s+%s\s*\))`,
	//	intervalFuncLiteral, timeInstanceFunc, intervalFuncLiteral, timeInstanceFunc, stringLiteral)
	//numberLiteral := `[0-9\.\-\+]+`
	//
	//// "test-1"."score"
	//columnName := fmt.Sprintf(`"%s"\."%s"`, word, word)
	//
	//// ("test-1"."score" + "test-2"."score")
	//arithmetic := fmt.Sprintf(`\((%s|\d+)\s+(\-|\+|/|\*)\s+(%s|\d+)\)`, columnName, columnName)
	//
	//aggregationFunc := fmt.Sprintf(`(?:(MAX\(\s*%s\s*\))|(MIN\(\s*%s\s*\))|(SUM\(\s*%s\s*\))|(AVG\(\s*%s\s*\))|(RANGE\(\s*%s\s*\))|(COUNT\(\s*\)))`,
	//	columnName, columnName, columnName, columnName, columnName)
	//
	//// ("test-1"."score" + "test-2"."score") AS "total"
	//// or "test-1"."score" AS "test-1"
	//// this fails because the group selectedCols is repeated.
	//colsOrArithmeticOrAggregateRegex := fmt.Sprintf(`(?:%s|%s|%s)\s+AS\s+"%s"`,
	//	aggregationFunc, columnName, arithmetic, word)
	//allColsRegex := `\*`
	//columnsRegex := fmt.Sprintf(`(?P<selectAllColumns>%s)|(?P<selectedColsOrArithmetic>%s(,\s+?%s)*)`,
	//	allColsRegex, colsOrArithmeticOrAggregateRegex, colsOrArithmeticOrAggregateRegex)
	//
	//// FROM "test-1"
	//tableRegex := fmt.Sprintf(`(?P<table>("%s"))`, word)
	//
	//// INNER/LEFT,RIGHT,FULL JOIN
	//// onConditionals
	//onConditionalsRegex := fmt.Sprintf(`\s+%s\s+=\s+(?:%s|%s|%s|%s)`,
	//	columnName, columnName, timeLiteral, stringLiteral, numberLiteral)
	//// joinsRegex := fmt.Sprintf(`(?P<join>((?:INNER|LEFT|RIGHT|FULL)\s+JOIN\s+"%s"\s+ON\s+%s(\s+AND\s+%s)*))`, word, onConditionalsRegex, onConditionalsRegex)
	//singleJoinRegex := fmt.Sprintf(`(?:INNER|LEFT|RIGHT|FULL)\s+JOIN\s+"%s"\s+ON%s(\s+AND%s)*`,
	//	word, onConditionalsRegex, onConditionalsRegex)
	//joinsRegex := fmt.Sprintf(`(?P<join>(\s+%s)*)`, singleJoinRegex)
	//
	//// GROUP BY
	//groupByRegex := fmt.Sprintf(`(?P<groupby>(\s+GROUP\s+BY\s+%s(?:\s*,\s*%s)*)*)`, columnName, columnName)
	//
	//// ORDER BY
	//colOrderBy := fmt.Sprintf(`%s\s+(DESC|ASC)`, columnName)
	//orderByRegex := fmt.Sprintf(`(?P<orderby>(\s+ORDER\s+BY\s+%s(?:\s*,\s*%s)*)*)`, colOrderBy, colOrderBy)
	//
	//// WHERE
	//comparator := `(?:(?:\<=)|(\<)|(?:\>=)|(\>)|(=))`
	//logicalOperator := `(?:AND|OR|NOT)`
	//colFilterRegex := fmt.Sprintf(`%s\s*%s\s*(?:%s|%s|%s|%s|%s)`,
	//	columnName, comparator, columnName, timeLiteral, stringLiteral, numberLiteral, arithmetic)
	//whereRegex := fmt.Sprintf(`(?P<where>(\s+WHERE\s*%s(?:\s*%s\s*%s)*)*)`, colFilterRegex, logicalOperator, colFilterRegex)
	//
	//regexString := fmt.Sprintf(
	//	`(?i)SELECT\s+%s\s+FROM\s+%s%s%s%s%s`, columnsRegex, tableRegex, whereRegex, joinsRegex, groupByRegex, orderByRegex,
	//)
	//
	//re := regexp.MustCompile(regexString)

	// FIXME: Looks like NOW() - INTERVAL('2 days') is not captured as one in WHERE
	//q := `SELECT
	//		MAX("student-details"."name") AS "name",
	//		SUM("test-1"."score") AS "test-1",
	//		AVG("test-2"."score") AS "test-2",
	//		"student-details"."name" AS "name",
	//		"test-1"."score" AS "test-1",
	//		"test-2"."score" AS "test-2",
	//		("test-1"."score" + "test-2"."score") AS "total"
	//	FROM "test-1"
	//	WHERE "test-2"."score" > NOW() - INTERVAL('2 days')
	//	INNER JOIN "test-2"
	//	ON "test-1"."student_id" = TO_TIMEZONE(NOW() + INTERVAL('2 years 30 days'), 'Europe/Paris')
	//		AND "test-1"."student_id" = "test-2"."student_id"
	//	INNER JOIN "student-details"
	//		ON "test-1"."student_id" = "student_details"."id"
	//	GROUP BY "student-details"."age", "student-details"."name"
	//	ORDER BY "student-details"."age" DESC, "student-details"."datetime" DESC
	//	`

	// q := `SELECT
	// 		MAX("student-details"."name") AS "name",
	// 		SUM("test-1"."score") AS "test-1",
	// 		AVG("test-2"."score") AS "test-2",
	// 		"student-details"."name" AS "name",
	// 		"test-1"."score" AS "test-1",
	// 		"test-2"."score" AS "test-2",
	// 		("test-1"."score" + "test-2"."score") AS "total"
	// 	  FROM "test-1"
	// 	  WHERE "test-2"."score" > NOW()
	// 	  INNER JOIN "test-2"
	// 	  ON "test-1"."student_id" = TO_TIMEZONE(NOW() + INTERVAL('2 years 30 days'), 'Europe/Paris')
	// 	  	AND "test-1"."student_id" = "test-2"."student_id"
	// 	  INNER JOIN "student-details"
	// 		ON "test-1"."student_id" = "student_details"."id"
	// 	  GROUP BY "student-details"."age", "student-details"."name"
	// 	  ORDER BY "student-details"."age" DESC, "student-details"."datetime" DESC
	// 	  `

	//submatches := re.FindStringSubmatch(q)
	//// fmt.Printf("length: %d, \n\nnames: %v", len(submatches), len(re.SubexpNames()))
	//result := make(map[string]string)
	//for i, name := range re.SubexpNames() {
	//	if i != 0 && name != "" {
	//		result[name] = submatches[i]
	//	}
	//}

	sqlFile := "example.sql"
	ql := &codegenerator.Eliql{}
	ql.RunFile(sqlFile)
}
