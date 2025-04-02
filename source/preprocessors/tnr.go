package preprocessors

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	db "example.com/source/database"
)

func GenerateTransitNodeRoutes(db *db.Postgres, concurrency int) {
	file, err := os.Create("tnr.csv")
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(file)
	w.Comma = 59 // ;
	if err := w.Write([]string{"id", "source", "target", "cost", "reverse_cost", "via"}); err != nil {
		panic(err)
	}

	id := -1

	defer file.Close()
	defer w.Flush()

	sem := make(chan struct{}, concurrency)
	driveways := getDriveways(db)
	exits := getExits(db)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, driveway := range driveways {
		for _, exit := range exits {
			sem <- struct{}{} // acquire slot
			wg.Add(1)

			go func(d, e int64) {
				defer wg.Done()
				defer func() { <-sem }() // release slot

				path := db.BdAstar(driveway, exit)
				mu.Lock()
				writeRow(w, id, driveway, exit, path)
				id--
				mu.Unlock()
			}(driveway, exit)
		}
	}

	wg.Wait()
}

func writeRow(w *csv.Writer, id int, source int64, target int64, path db.Path) {
	cost := 0.0
	if len(path.Sequences) > 0 {
		cost = path.Sequences[len(path.Sequences)-1].AggCost
	}

	row := []string{
		fmt.Sprintf("%d", id),
		fmt.Sprintf("%d", source),
		fmt.Sprintf("%d", target),
		fmt.Sprintf("%f", cost),
		"-1",           // reverse cost
		path.ToArray(), // via
	}

	if err := w.Write(row); err != nil {
		panic(err)
	}
}

func getDriveways(db *db.Postgres) []int64 {
	rows, err := db.Query(`
		WITH Junctions AS (
		    SELECT f_jnctid
		    FROM nw
		    WHERE frc = 0 OR routenum = 'R0'
		)
		SELECT
		    CASE
		         WHEN nl.oneway = 'FT' THEN nw.t_jnctid
		         WHEN nl.oneway = 'TF' THEN nw.f_jnctid
		         ELSE NULL
		    END AS junction_id
		FROM nw
		JOIN nl ON nw.id = nl.id
		WHERE nw.frc <> 0
		  AND (nw.f_jnctid IN (SELECT f_jnctid FROM Junctions) OR nw.t_jnctid IN (SELECT f_jnctid FROM Junctions))
		  AND (
		    (nl.oneway = 'FT' AND nw.t_jnctid IN (SELECT f_jnctid FROM Junctions))
		    OR
		    (nl.oneway = 'TF' AND nw.f_jnctid IN (SELECT f_jnctid FROM Junctions))
		  );
	`)
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
	rows, err := db.Query(`
		WITH Junctions AS (
		    SELECT f_jnctid
		    FROM nw
		    WHERE frc = 0 OR routenum = 'R0'
		)
		SELECT
		  CASE
		    WHEN nl.oneway = 'FT' THEN nw.f_jnctid
		    WHEN nl.oneway = 'TF' THEN nw.t_jnctid
		  END AS junction_id
		FROM nw
		JOIN nl ON nw.id = nl.id
		WHERE nw.frc <> 0
		  AND (
		    nw.f_jnctid IN (SELECT f_jnctid FROM Junctions)
		    OR
		    nw.t_jnctid IN (SELECT f_jnctid FROM Junctions)
		  )
		  AND (
		    (nl.oneway = 'FT' AND nw.f_jnctid IN (SELECT f_jnctid FROM Junctions))
		    OR
		    (nl.oneway = 'TF' AND nw.t_jnctid IN (SELECT f_jnctid FROM Junctions))
		  );
	`)
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
