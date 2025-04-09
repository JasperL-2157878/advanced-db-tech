package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

var shortcutsFile, graphFile, accessFile *os.File
var shortcuts, graph, accesses *csv.Writer

type Edge struct {
	ID, Source, Target                int64
	Cost, ReverseCost, X1, Y1, X2, Y2 float64
}

type AccessNode struct {
	Junction int64
	Cost     float64
}

var junctions []int64
var edges map[int64]Edge = getEdges()
var driveways []int64 = getDriveways()
var exits []int64 = getExits()

func open() {
	shortcutsFile, _ = os.Create("tnr_shortcuts.csv")
	shortcuts = csv.NewWriter(shortcutsFile)
	shortcuts.Comma = ';'
	shortcuts.Write([]string{"id", "source", "target", "cost", "reverse_cost", "nodes", "edges"})

	graphFile, _ = os.Create("tnr_graph.csv")
	graph = csv.NewWriter(graphFile)
	graph.Comma = ';'
	graph.Write([]string{"id", "source", "target", "cost", "reverse_cost", "x1", "y1", "x2", "y2"})

	accessFile, _ = os.Create("tnr_access.csv")
	accesses = csv.NewWriter(accessFile)
	accesses.Comma = ';'
	accesses.Write([]string{"junction", "access_source_nodes", "access_source_costs", "access_target_nodes", "access_target_costs"})
}

func close() {
	for _, edge := range edges {
		graph.Write([]string{
			fmt.Sprintf("%d", edge.ID),
			fmt.Sprintf("%d", edge.Source),
			fmt.Sprintf("%d", edge.Target),
			fmt.Sprintf("%f", edge.Cost),
			fmt.Sprintf("%f", edge.ReverseCost),
			fmt.Sprintf("%f", edge.X1),
			fmt.Sprintf("%f", edge.Y1),
			fmt.Sprintf("%f", edge.X2),
			fmt.Sprintf("%f", edge.Y2),
		})
	}

	shortcuts.Flush()
	shortcutsFile.Close()

	graph.Flush()
	graphFile.Close()

	accesses.Flush()
	accessFile.Close()
}

func tnr() {
	open()

	sem := make(chan struct{}, 8)
	n := len(driveways)
	m := len(exits)

	id := int64(-1)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, driveway := range driveways {
		for j, exit := range exits {
			fmt.Println(i*n+j, "/", n*m)

			sem <- struct{}{} // acquire slot
			wg.Add(1)

			go func(id, d, e int64) {
				defer wg.Done()
				defer func() { <-sem }() // release slot

				if d == e {
					return
				}

				path := Graphs.ChDijkstra(driveway, exit)
				if path.Cost <= 0 || len(path.Edges) <= 1 {
					return
				}

				path.Nodes = path.Nodes[1 : len(path.Nodes)-1]
				path.Edges = path.Edges[1 : len(path.Edges)-1]

				mu.Lock()
				shortcuts.Write([]string{
					fmt.Sprintf("%d", id),
					fmt.Sprintf("%d", d),
					fmt.Sprintf("%d", e),
					fmt.Sprintf("%f", path.Cost),
					"-1",
					intArrayToString(path.Nodes),
					intArrayToString(path.Edges),
				})

				for _, edge := range path.Edges {
					delete(edges, edge)
				}

				id--
				mu.Unlock()
			}(id, driveway, exit)
		}
	}
	wg.Wait()

	/*junctions = getJunctions()
	for _, junction := range junctions {
		sem <- struct{}{} // acquire slot
		wg.Add(1)

		go func(junction int64) {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			var accessSources, accessTargets []AccessNode

			for _, driveway := range driveways {
				path := Graphs.ChDijkstra(junction, driveway)
				accessSources = append(accessSources, AccessNode{
					Junction: driveway,
					Cost:     path.Cost,
				})
			}

			for _, exit := range exits {
				path := Graphs.ChDijkstra(junction, exit)
				accessSources = append(accessSources, AccessNode{
					Junction: exit,
					Cost:     path.Cost,
				})
			}

			sort.Slice(accessSources, func(i, j int) bool {
				return accessSources[i].Cost < accessSources[j].Cost
			})

			sort.Slice(accessTargets, func(i, j int) bool {
				return accessTargets[i].Cost < accessTargets[j].Cost
			})

			mu.Lock()
			shortcuts.Write([]string{
				fmt.Sprintf("%d", junction),
				accessArrayToString(accessSources[:6]),
				accessArrayToString(accessTargets[:6]),
			})
			mu.Unlock()
		}(junction)
	}
	wg.Wait()*/

	close()
}

func getJunctions() []int64 {
	rows, err := Db.Query(`
		SELECT DISTINCT id
		FROM (
			SELECT source AS id FROM tnr_graph
			UNION
			SELECT target AS id FROM tnr_graph
		)
	`)
	if err != nil {
		panic(err)
	}

	var junctions []int64 = make([]int64, 0)
	for rows.Next() {
		var junction int64
		err = rows.Scan(&junction)
		if err != nil {
			panic(err)
		}
		junctions = append(junctions, junction)
	}

	return junctions
}

func getEdges() map[int64]Edge {
	rows, err := Db.Query(`SELECT * FROM base_graph`)
	if err != nil {
		panic(err)
	}

	var edges map[int64]Edge = make(map[int64]Edge)
	for rows.Next() {
		var edge Edge
		err = rows.Scan(&edge.ID, &edge.Source, &edge.Target, &edge.Cost, &edge.ReverseCost, &edge.X1, &edge.Y1, &edge.X2, &edge.Y2)
		if err != nil {
			panic(err)
		}
		edges[edge.ID] = edge
	}

	return edges
}

func getDriveways() []int64 {
	rows, err := Db.Query(`
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

func getExits() []int64 {
	rows, err := Db.Query(`
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

func intArrayToString(items []int64) string {
	str := "{"

	for i, item := range items {
		str += fmt.Sprintf("%d", item)
		if i < len(items)-1 {
			str += ","
		}
	}

	return str + "}"
}

func accessArrayToString(items []AccessNode) string {
	nodes := "{"
	costs := "{"

	for i, item := range items {
		nodes += fmt.Sprintf("%d", item.Junction)
		costs += fmt.Sprintf("%f", item.Cost)
		if i < len(items)-1 {
			nodes += ","
			costs += ","
		}
	}

	return nodes + "};" + costs + "}"
}
