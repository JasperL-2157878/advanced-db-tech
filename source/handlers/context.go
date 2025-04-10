package handlers

import (
	"net/http"

	db "example.com/source/database"
	"example.com/source/graph"
)

type Context struct {
	req   *http.Request
	res   http.ResponseWriter
	db    *db.Postgres
	graph *graph.Graph
}

func NewContext(r *http.Request, w http.ResponseWriter, db *db.Postgres, graph *graph.Graph) Context {
	return Context{
		req:   r,
		res:   w,
		db:    db,
		graph: graph,
	}
}

func (c *Context) Param(param string) string {
	return c.req.URL.Query().Get(param)
}
