package database

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jWrapper struct {
	driver neo4j.DriverWithContext
	ctx    context.Context
}

func newNeo4j() *Neo4jWrapper {
	host := "host"
	user := "user"
	pass := "pass"

	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(host, neo4j.BasicAuth(user, pass, ""))

	if err != nil {
		panic(err)
	}

	return &Neo4jWrapper{driver: driver, ctx: ctx}
}

func (conn *Neo4jWrapper) Close() {
	conn.driver.Close(conn.ctx)
}

func (conn *Neo4jWrapper) Exec(sql string, params map[string]any) Summary {
	result, err := neo4j.ExecuteQuery(conn.ctx, conn.driver, sql, params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)

	if err != nil {
		panic(err)
	}

	summary := result.Summary.Counters()
	changes := summary.LabelsAdded() +
		summary.LabelsRemoved() +
		summary.PropertiesSet() +
		summary.RelationshipsCreated() +
		summary.RelationshipsDeleted() +
		summary.NodesCreated() +
		summary.NodesDeleted()

	return Summary{
		LastInsertId: -1,
		Changes:      int64(changes),
	}
}

func (conn *Neo4jWrapper) Query(sql string, params map[string]any) []Record {
	results, err := neo4j.ExecuteQuery(conn.ctx, conn.driver, sql, params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)

	if err != nil {
		panic(err)
	}

	var records []Record
	for _, r := range results.Records {
		records = append(records, r.AsMap())
	}

	return records
}

func (conn *Neo4jWrapper) QuerySingle(sql string, params map[string]any) Record {
	results, err := neo4j.ExecuteQuery(conn.ctx, conn.driver, sql, params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)

	if err != nil {
		panic(err)
	}

	if len(results.Records) == 0 {
		return make(Record)
	} else {
		return results.Records[0].AsMap()
	}
}
