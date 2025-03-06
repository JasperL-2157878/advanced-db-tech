package handlers

import (
	"net/http"
)

func HandleIndex(res http.ResponseWriter, req *http.Request) {
	r1 := Db.Pgsql.Query("SELECT * from actors")
	r2 := Db.Pgsql.QuerySingle("SELECT * from actors")

	res.Write(JSON(map[string]any{
		"Hello": "World",
		"r1":    r1,
		"r2":    r2,
	}))
}
