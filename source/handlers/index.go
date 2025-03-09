package handlers

import (
	"net/http"
)

func HandleIndex(res http.ResponseWriter, req *http.Request) {
	res.Write(Db.Example())
}
