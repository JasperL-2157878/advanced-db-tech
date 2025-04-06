package handlers

import (
	"net/http"

	db "example.com/source/database"
	"example.com/source/graphs"
)

type Context struct {
	req    *http.Request
	res    http.ResponseWriter
	db     *db.Postgres
	graphs *graphs.Graphs
}

func NewContext(r *http.Request, w http.ResponseWriter, db *db.Postgres, graphs *graphs.Graphs) Context {
	return Context{
		req:    r,
		res:    w,
		db:     db,
		graphs: graphs,
	}
}

func (c *Context) Param(param string) string {
	return c.req.URL.Query().Get(param)
}
