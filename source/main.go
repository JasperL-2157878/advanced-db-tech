package main

import (
	"log"
	"net/http"

	"example.com/source/handlers"
)

func JSON(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		handler(w, r)
	}
}

func main() {
	defer handlers.Db.Close()

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/api", JSON(handlers.HandleIndex))

	log.Println("[INFO] Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
