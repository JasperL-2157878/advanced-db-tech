package db

import (
	"log"
	"sync"
)

func (pg *PostgresConnection) GenerateTNR(concurrency int) {
	sem := make(chan struct{}, concurrency)
	driveAways := pg.getDriveAways()
	exits := pg.getExits()

	log.Println("Drive aways amount: %d", len(driveAways))
	log.Println("Exits amount: %d", len(exits))

	var wg sync.WaitGroup
	for _, driveAway := range driveAways {
		for _, exit := range exits {
			sem <- struct{}{} // acquire slot
			wg.Add(1)

			go func(d, e int) {
				defer wg.Done()
				defer func() { <-sem }() // release slot
				pg.generateTNRPath(d, e)
			}(driveAway, exit)
		}
	}
	wg.Wait()
}

func (pg *PostgresConnection) generateTNRPath(driveAway int, exit int) {
	route := pg.Route(driveAway, exit)

	if len(route.Features) == 0 {
		log.Println("No route found between %d and %d", driveAway, exit)
		return
	}

	currentRouteNum := route.Features[0].Properties.RouteNum

	for _, feature := range route.Features {
		if feature.Properties.RouteNum != currentRouteNum {
			return
		}
	}

	log.Println("Route found between %d and %d and total cost is %f", driveAway, exit, route.TotalCost)
}

func (pg *PostgresConnection) getDriveAways() []int {
	rows, err := pg.conn.Query(`
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

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}

	return ids
}

func (pg *PostgresConnection) getExits() []int {
	rows, err := pg.conn.Query(`
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

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}

	return ids
}
