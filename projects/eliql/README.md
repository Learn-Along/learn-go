# EliQL

The Query language for stream processing data in golang e.g.

```SQL
    SELECT
        MAX("day_of_year") AS "day_of_year",
        MIN("day_of_week") AS "day_of_week",
        AVG("unit price") AS "unit price",
        SUM("amount") AS "amount",
        PRODUCT(
          DIFFERENCE("amount", "cost"),
          AVG("unit price")
        ) AS "net profit",
        NOW() AS "created_at"
    FROM "foo" LEFT JOIN "bar" ON "foo"."datetime" = "bar"."datetime"
        GROUP BY "datetime", "timezone"
        ORDER BY "datetime" DESC, "timezone" ASC
```

## Design

The syntax that is very similar to SQL.
The syntax follows a similar structure as SQL and so a [short tutorial about SQL](https://www.w3schools.com/sql/) will get you very far.

However there are a few marked differences. They include:

### Only SELECT is allowed

- Only the `SELECT` statement is allowed. No mutation is expected in this processing thus the other
  statements like `CREATE TABLE`, `INSERT` etc are useless.
- Quotation marks are a must on each parameter except for the special case "\*" which selects all columns
- `"table_name".` prefix is a must for queries that contain `JOIN`. In the background, the prefix is added for all queries but this is derived from the `FROM` statement.
- Arithmetic computations are always of form `(...) AS "..."`
- Function calls are always of the form `function_name(...) AS "..."`
- All columns must be of any of the following forms:
  - `"column_name" AS "..."`
  - `"table"."column_name" AS "..."`
  - `(...) AS "..."` e.g. `("foo" - "goo") AS "bar"`, `("results"."foo" - "results"."goo") AS "bar"`
  - `function_name(...)` e.g. `SUM("foo") AS "foo"`, `SUM("results"."foo") AS "foo"`
- In order to esure that the processed records each have understandable field names, each _'column'_ in the select statement must have an `AS ...` suffix. e.g.

  ```SQL
  SELECT "amount" AS "amount"
  FROM "foo"
  ```

- However, the `AS ...` suffix should never be put on the `*` selector e.g. the EliQL below is perfectly valid

  ```SQL
  SELECT * FROM "foo"
  ```

### _tables_

- the _'tables'_ in this case are "message types" as labelled in the collection.
- these message types have to be surrounded by `""` (double quotation marks). This is to cater for message types that have spaces within them e.g. "foo and bar"

### _columns_

- the _'columns'_ in this case are the individual properties expected on each record in a given message type.
- these properties have to be surrounded by `""` (double quotation marks). This is to cater for properties that have spaces within them e.g. "amount and cost"
- In some instances where a property name might be confused for another in another message, the
  property is denoted with the message type as a prefix plus a dot e.g. "foo"."amount".
  This kind of notation is mandatory when putting conditions for a `JOIN` operation e.g.

  ```SQL
  SELECT * FROM "foo" INNER JOIN "bar" ON "foo"."datetime" = "bar"."datetime"
  ```

### Artithmetic

- `+`, `-`, `*`, `/` are used.
- The computation is always surrounded by parentheses `()`
- These arithmetic functions can easily be nested e.g.

  ```SQL
  SELECT
    (("amount" + "price") - "total amount" - ("day_of_year" / 365) * 100)) AS "some property"
  FROM "foo"
  ```

### Aggregation

- Aggregation is mainly done through `GROUP BY` and some in-built functions as shown in the table below

| Function     | parameters | Use                                                                                           |
| ------------ | ---------- | --------------------------------------------------------------------------------------------- |
| MIN(field)   | field      | Returns the minimum value for the `field` in each grouping                                    |
| MAX(field)   | field      | Returns the maximum value for the `field` in each grouping                                    |
| AVG(field)   | field      | Returns the average value of the field in each grouping                                       |
| RANGE(field) | field      | Returns the difference between the biggest and smallest value for the field for each grouping |
| SUM(field)   | field      | Returns the sum of the field in each grouping                                                 |
| COUNT()      |            | Returns the number of records in each grouping                                                |

An example of aggregation is as follows:

```SQL
SELECT
  AVG("price") AS "price",
  SUM("amount") AS "amount"
FROM "foo"
  GROUP BY "foo"."datetime"
```

### Sorting

- Sorting is done using the usual `ORDER BY` with the fields followed by either `DESC` for descending order sort or `ASC` for ascending order sort. An example is:

  ```SQL
  SELECT *
  FROM "foo"
    ORDER BY "foo"."amount" DESC, "foo"."datetime" DESC
  ```

### Filtering

- Filtering is done using the usual `WHERE` clause togther with comparator signs and logic operators as shown in the table below

| Operator | Use                                                                          |
| -------- | ---------------------------------------------------------------------------- |
| >        | Greater than                                                                 |
| >=       | Greater or equal to                                                          |
| <        | Less than                                                                    |
| <=       | Less or equal to                                                             |
| =        | Equal to                                                                     |
| AND      | the operands on either side of `AND` must be true to produce true            |
| OR       | Any of the operand on any side of `OR` will cause the result to be true also |
| NOT      | the result is true if the oerand one right is false, and vice versa          |

An example of filtering is:

```SQL
SELECT *
FROM "foo"
  WHERE "foo"."amount" > 0
```

### Time and intervals

- In order to have queries that change given on the current timestamp, we have a number of time and interval functions as listed in the table below.

| Function                                     | parameters                                                                                                                                                                                                             | Use                                                                                  |
| -------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| NOW()                                        |                                                                                                                                                                                                                        | Returns the timestamp now based on the timezone of the system the app is running on. |
| TO_TIMEZONE(timestamp/date, timezone_string) | first parameter is the timestamp or date e.g, from NOW() or TODAY() while the second parameter is the timezone string e.g. 'Europe/Paris'                                                                              | Converts the given timestamp to the given timezone                                   |
| TODAY()                                      |                                                                                                                                                                                                                        | Returns the date today based on the timezone of the system the app is running on.    |
| INTERVAL(string)                             | string representing the interval in form of `[number] years [number] days [number] hours [number] minutes [number] seconds [number] milliseconds` where all are optional. Valid example usage is `INTERVAL('2 hours')` | Returns a time duration basing on the argument string passed to it.                  |

An example of a time-dependent query is:

```SQL
SELECT
  (TO_TIMEZONE(NOW(), 'Europe/Paris') - INTERVAL('2 days 78 hours')) AS `created at`
FROM `foo`
  WHERE `amount` > 0
```

### Joins

EliQL supports `JOIN`, `LEFT JOIN`, `RIGHT JOIN`, `FULL JOIN`. More details on these can be found in the tutorial at [w3schools.com](https://www.w3schools.com/sql/sql_join.asp)

### Unions

EliQL also supports combining two queries in one using the `UNION` (no duplicates) and `UNION ALL` (duplicates allowed)

### Other In-built Functions

There are other in built functions that are neither for aggregation nor for time intervals. They include:

| Function    | parameters      | Use                                                                                       |
| ----------- | --------------- | ----------------------------------------------------------------------------------------- |
| CONCAT(...) | fields, strings | Concatenates the values in the fields with each other and with any string literals passed |

### Comments

Comments are prefixed by `--`

e.g.

```SQL
-- This is a comment
SELECT * FROM "db"
```

## Syntax

Here is the current syntax grammar to be followed as shown in [syntax.txt](./syntax.txt).

```
expression              -> unionExpr | selectExpr;
unionExpr               -> selectExpr (union selectExpr)+;
selectExpr              -> "SELECT" columnExpr ("," columnExpr)* "FROM" TABLE join* where? groupBy? orderBy?;
union                   -> "UNION" "ALL"?;
columnExpr              -> (arithmetic | COLUMN | FUNCTION) "AS" NAME;
join                    -> ("LEFT" | "RIGHT" | "FULL")? "JOIN" TABLE
                            "ON" columnEqualToColumn ("AND" columnEqualToColumn)*;
where                   -> "WHERE" comparison (logicalOperator comparison)*;
groupBy                 -> "GROUP" "BY" COLUMN ("," COLUMN)*;
orderBy                 -> "ORDER" "BY" columnOrder ("," columnOrder)*;
arithmetic              -> "(" COLUMN (arithmeticOperator (COLUMN | NUMBER | STRING))+ ")";
columnEqualToColumn     -> COLUMN "=" COLUMN
comparison              -> "NOT"? COLUMN comparator (COLUMN | NUMBER | STRING);
logicalOperator         -> "AND" | "OR";
comparator              -> ">" | ">=" | "=" | "<" | "<=";
columnOrder             -> COLUMN ("ASC" | "DESC")
arithmeticOperator      -> "/" | "-" | "+" | "*"
```

## How to Use

- Install the package

```sh
go get github.com/Learn-Along/learn-go/projects/eliql
```

- Create a stream from say a websocket connection or a REST API e.g.

```go
package main

import (
  "github.com/Learn-Along/learn-go/projects/eliql"
  "time"
  "log"
)

func main() {
  bufferSize := 5
  // create a collection which caches stream messages in a buffer;
  // bufferSize is the maximum number of messages in buffer per stream
  collection := eliql.NewCollection(bufferSize)

  // Add a stream of messages received from a websocket
  err := collection.AddWebsocketStream("test-1", "ws://example.com/results/test-1")
  if err != nil {
    log.Fatal("error test-1: ", err)
  }

  // Add a stream of messages received from a websocket
  err := collection.AddWebsocketStream("test-2", "ws://example.com/results/test-2")
  if err != nil {
    log.Fatal("error test-2: ", err)
  }

  // Add a stream of messages received from a REST API being queried at the given interval
  err := collection.AddRestAPIStream("student-details", "https://example.com/students/details", 30 * time.Second)

  // Create a new stream, a result of processing the above streams
  totalsStream, err := collection.Query(
    `SELECT
      "student-details"."name" AS "name",
      "test-1"."score" AS "test-1",
      "test-2"."score" AS "test-2",
      ("test-1"."score" + "test-2"."score") AS "total"
    FROM "test-1"
    INNER JOIN "test-2"
      ON "test-1"."student_id" = "test-2"."student_id"
    INNER JOIN "student-details"
      ON "test-1"."student_id" = "student_details"."id"
    `
  )
  if err != nil {
    log.Fatal("error test-2: ", err)
  }

  for record := range totalsStream.C {
    // Print the details for each student as a map {"name": "...", "test-1": ..., "test-2": ..., "total": ...}
    fmt.Printf("%v", record)
  }

}
```

## License

Copyright (c) 2021 [Martin Ahindura](https://github.com/Tinitto) Licensed under the [MIT License](./LICENSE)
