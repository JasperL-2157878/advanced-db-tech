package database

type DatabasePool struct {
	Pgsql *PostgresWrapper
	Neo4j *Neo4jWrapper
}

func Connect() *DatabasePool {
	return &DatabasePool{
		Pgsql: newPostgreSQL(),
		Neo4j: newNeo4j(),
	}
}

func (db *DatabasePool) Close() {
	db.Pgsql.Close()
	db.Neo4j.Close()
}

type Record map[string]any
type Summary struct {
	LastInsertId int64
	Changes      int64
}

type DatabaseWrapper interface {
	Exec(sql string, params ...map[string]any) Summary
	Query(sql string, params ...map[string]any) []Record
	QuerySingle(sql string, params ...map[string]any) Record
	Close()
}
