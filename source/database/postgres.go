package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PostgresWrapper struct {
	driver *sql.DB
}

func newPostgreSQL() *PostgresWrapper {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	log.Println("[INFO] Connecting to Postgres ...")
	connstr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", user, password, dbname)
	driver, err := sql.Open("postgres", connstr)
	if err != nil {
		panic(err)
	}

	err = driver.Ping()
	if err != nil {
		panic(err)
	}

	return &PostgresWrapper{
		driver: driver,
	}
}

func (conn *PostgresWrapper) Close() {
	conn.driver.Close()
}

func (conn *PostgresWrapper) Exec(sql string, params ...map[string]any) Summary {
	result, err := conn.driver.Exec(sql, conn.toArgs(params)...)
	if err != nil {
		panic(err)
	}

	var summary Summary
	summary.LastInsertId, _ = result.LastInsertId()
	summary.Changes, _ = result.RowsAffected()

	return summary
}

func (conn *PostgresWrapper) Query(sql string, params ...map[string]any) []Record {
	rows, err := conn.driver.Query(sql, conn.toArgs(params)...)
	if err != nil {
		panic(err)
	}

	var records []Record
	for rows.Next() {
		records = append(records, conn.toMap(rows))
	}

	return records
}

func (conn *PostgresWrapper) QuerySingle(sql string, params ...map[string]any) Record {
	rows, err := conn.driver.Query(sql, conn.toArgs(params)...)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		return conn.toMap(rows)
	} else {
		return nil
	}
}

func (conn *PostgresWrapper) toMap(rows *sql.Rows) map[string]any {
	fields, _ := rows.Columns()

	scans := make([]interface{}, len(fields))
	record := make(Record)

	for i := range scans {
		scans[i] = &scans[i]
	}

	rows.Scan(scans...)
	for i, v := range scans {
		var value = ""
		if v != nil {
			value = fmt.Sprintf("%s", v)
		}
		record[fields[i]] = value
	}

	return record
}

func (conn *PostgresWrapper) toArgs(params []map[string]any) []any {
	if len(params) == 0 {
		return make([]any, 0)
	}

	data := make([]any, len(params[0]))
	for _, val := range params[0] {
		data = append(data, val)
	}

	return data
}
