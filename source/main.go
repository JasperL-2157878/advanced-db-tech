package main

import (
	"log"
	"net/http"

	"example.com/source/handlers"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/api", handlers.HandleIndex)

	log.Println("[INFO] Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
