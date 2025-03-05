package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgreSQLWrapper struct {
	driver *sql.DB
}

func newPostgreSQL() *PostgreSQLWrapper {
	driver, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		panic(err)
	}

	return &PostgreSQLWrapper{
		driver: driver,
	}
}

func (conn *PostgreSQLWrapper) Close() {
	conn.driver.Close()
}

func (conn *PostgreSQLWrapper) Exec(sql string, params map[string]any) Summary {
	result, err := conn.driver.Exec(sql, params)
	if err != nil {
		panic(err)
	}

	var summary Summary
	summary.LastInsertId, _ = result.LastInsertId()
	summary.Changes, _ = result.RowsAffected()

	return summary
}

func (conn *PostgreSQLWrapper) Query(sql string, params map[string]any) []Record {
	rows, err := conn.driver.Query(sql, params)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		rows.Scan(record)
		records = append(records, record)
	}

	return records
}

func (conn *PostgreSQLWrapper) QuerySingle(sql string, params map[string]any) Record {
	row := conn.driver.QueryRow(sql, params)

	var record Record
	row.Scan(record)

	return record
}
