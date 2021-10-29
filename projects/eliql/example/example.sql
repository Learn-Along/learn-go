SELECT
    MAX("student-details"."name") AS "name",
    SUM("test-1"."score") AS "test-1",
    AVG("test-2"."score") AS "test-2",
    "student-details"."name" AS "name",
    "test-1"."score" AS "test-1",
    "test-2"."score" AS "test-2",
    ("test-1"."score" + "test-2"."score") AS "total"
FROM "test-1"
WHERE "test-2"."score" > NOW() - INTERVAL('2 days')
    INNER JOIN "test-2"
ON "test-1"."student_id" = TO_TIMEZONE(NOW() + INTERVAL('2 years 30 days'), 'Europe/Paris')
    AND "test-1"."student_id" = "test-2"."student_id"
    INNER JOIN "student-details"
    ON "test-1"."student_id" = "student_details"."id"
--  This is a comment
GROUP BY "student-details"."age", "student-details"."name"
ORDER BY "student-details"."age" DESC, "student-details"."datetime" DESC;