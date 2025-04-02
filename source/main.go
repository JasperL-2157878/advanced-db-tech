package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/source/handlers"
	"example.com/source/preprocessors"
)

func preprocess() {
	fmt.Println("Preprocessing (could take multiple hours) ...")
	preprocessors.GenerateTransitNodeRoutes(handlers.Db, 8)
	preprocessors.GenerateContractionHierarchies(handlers.Db)
}

func JSON(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		handler(w, r)
	}
}

func main() {
	defer handlers.Db.Close()

	preprocess()
	os.Exit(0)

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/api/v1/geocode", JSON(handlers.HandleGeocode))
	http.HandleFunc("/api/v1/route", JSON(handlers.HandleDijkstra))

	http.HandleFunc("/api/v1/route/dijkstra", JSON(handlers.HandleDijkstra))
	http.HandleFunc("/api/v1/route/astar", JSON(handlers.HandleAstar))
	http.HandleFunc("/api/v1/route/bddijkstra", JSON(handlers.HandleBdDijkstra))
	http.HandleFunc("/api/v1/route/bdastar", JSON(handlers.HandleBdAstar))

	//http.HandleFunc("/api/v1/route/raw", JSON(handlers.HandleBdDijkstra))
	//http.HandleFunc("/api/v1/route/tnr", JSON(handlers.HandleTnr))
	//http.HandleFunc("/api/v1/route/ch", JSON(handlers.HandleCh))
	//http.HandleFunc("/api/v1/route/tnrch", JSON(handlers.HandleTnrCh))

	log.Println("[INFO] Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
