package handlers

import (
	"net/http"
)

func HandleIndex(res http.ResponseWriter, req *http.Request) {
	Db.Pgsql.Close()

	res.Header().Add("content-type", "application/json")
	res.Write(JSON(map[string]any{
		"Hello": "World",
	}))
}
