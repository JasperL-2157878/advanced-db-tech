package main

import (
	"log"
	"net/http"

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
	}
}

/* func foo() {
	file, err := os.Create("tnr.csv")
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(file)
	w.Comma = 59 // ;
	if err := w.Write([]string{"id", "source", "target", "cost", "reverse_cost", "nodes", "edges"}); err != nil {
		panic(err)
	}

	defer file.Close()
	defer w.Flush()

	sem := make(chan struct{}, 8)
	driveways := getDriveways(Db)
	exits := getExits(Db)

	g := graphs.LoadBaseGraph(Db)

	var wg sync.WaitGroup
	var mu sync.Mutex

	n := len(driveways)
	id := int64(-1)

	for i := 0; i < n; i++ {
		sem <- struct{}{} // acquire slot
		wg.Add(1)

		driveway := driveways[i]
		exit := exits[i]

		go func(d, e int64) {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			path := g.BdDijkstra(driveway, exit)
			mu.Lock()
			writeRow(w, id, driveway, exit, path)
			id--
			mu.Unlock()
		}(driveway, exit)
	}

	wg.Wait()
}

func writeRow(w *csv.Writer, id int64, source int64, target int64, path *types.Path) {
	row := []string{
		fmt.Sprintf("%d", id),
		fmt.Sprintf("%d", source),
		fmt.Sprintf("%d", target),
		fmt.Sprintf("%f", path.Cost),
		"-1",
		path.NodesToArray(),
		path.EdgesToArray(),
	}

	if err := w.Write(row); err != nil {
		panic(err)
	}
}

func getDriveways(db *db.Postgres) []int64 {
	rows, err := db.Query(`SELECT source FROM tnr_old ORDER BY id`)
	if err != nil {
		panic(err)
	}

	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}

	return ids
}

func getExits(db *db.Postgres) []int64 {
	rows, err := db.Query(`SELECT target FROM tnr_old ORDER BY id`)
	if err != nil {
		panic(err)
	}

	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}

	return ids
} */

func main() {
	defer Db.Close()

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/api/v1/geocode", di(handlers.HandleGeocode))
	http.HandleFunc("/api/v1/route", di(handlers.HandleAlgDijkstra))

	http.HandleFunc("/api/v1/route/alg/dijkstra", di(handlers.HandleAlgDijkstra))
	http.HandleFunc("/api/v1/route/alg/astar", di(handlers.HandleAlgAstar))
	http.HandleFunc("/api/v1/route/alg/bddijkstra", di(handlers.HandleAlgBdDijkstra))
	http.HandleFunc("/api/v1/route/alg/bdastar", di(handlers.HandleAlgBdAstar))

	http.HandleFunc("/api/v1/route/opt/none", di(handlers.HandleOptNone))
	http.HandleFunc("/api/v1/route/opt/tnr", di(handlers.HandleOptTnr))
	//http.HandleFunc("/api/v1/route/opt/ch", di(handlers.HandleOptCh))
	//http.HandleFunc("/api/v1/route/opt/full", di(handlers.HandleOptFull))

	log.Println("[INFO] Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
