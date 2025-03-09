package handlers

import (
	"net/http"
)

func HandleIndex(res http.ResponseWriter, req *http.Request) {
	res.Write(JSON(map[string]any{
		"Hello": "World",
	}))
}
