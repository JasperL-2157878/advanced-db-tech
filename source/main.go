package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	db "example.com/source/database"
	"example.com/source/graphs"
	"example.com/source/handlers"
)

var Db *db.Postgres = db.NewPostgres()
var Graphs *graphs.Graphs = graphs.LoadGraphs(Db)

// Injects dependencies into Context-handlers and translates them to HandleFunc-handlers
func di(f func(ctx handlers.Context)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		f(handlers.NewContext(r, w, Db, Graphs))
		runtime.GC() // somewhere memory leak :(
	}
}

func main() {
	defer Db.Close()

	tnr()
	os.Exit(0)

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/api/v1/geocode", di(handlers.HandleGeocode))
	http.HandleFunc("/api/v1/route", di(handlers.HandleAlgDijkstra))

	http.HandleFunc("/api/v1/route/alg/dijkstra", di(handlers.HandleAlgDijkstra))
	http.HandleFunc("/api/v1/route/alg/astar", di(handlers.HandleAlgAstar))
	http.HandleFunc("/api/v1/route/alg/bddijkstra", di(handlers.HandleAlgBdDijkstra))
	http.HandleFunc("/api/v1/route/alg/bdastar", di(handlers.HandleAlgBdAstar))

	http.HandleFunc("/api/v1/route/opt/none", di(handlers.HandleOptNone))
	http.HandleFunc("/api/v1/route/opt/tnr", di(handlers.HandleOptTnr))
	http.HandleFunc("/api/v1/route/opt/ch", di(handlers.HandleOptCh))
	http.HandleFunc("/api/v1/route/opt/chtnr", di(handlers.HandleOptChTnr))

	log.Println("[INFO] Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
