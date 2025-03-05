package main

import (
	"net/http"

	"example.com/source/handlers"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/api", handlers.HandleIndex)

	http.ListenAndServe(":8080", nil)
}
