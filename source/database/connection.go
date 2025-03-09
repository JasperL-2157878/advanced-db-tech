package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	conn *sql.DB
}

func NewPostgresConnection() *PostgresConnection {
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

	return &PostgresConnection{
		driver,
	}
}

func (pg *PostgresConnection) Close() {
	pg.conn.Close()
}
