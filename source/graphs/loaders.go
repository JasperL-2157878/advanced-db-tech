package graphs

import (
	db "example.com/source/database"
	"github.com/LdDl/ch"
	"github.com/lib/pq"
)

func loadEdges(g *Graphs, db *db.Postgres) {
	rows, err := db.Query(`SELECT id, source, target FROM base_graph`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, source, target int64

		err := rows.Scan(&id, &source, &target)
		if err != nil {
			panic(err)
		}

		g.edges[[2]int64{source, target}] = id
		g.edges[[2]int64{target, source}] = id
	}

	rows.Close()
}

func loadTnrShortcuts(g *Graphs, db *db.Postgres) {
	rows, err := db.Query(`SELECT source, target, cost, nodes, edges FROM tnr_shortcuts`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var source, target int64
		var cost float64
		var nodes, edges pq.Int64Array

		err := rows.Scan(&source, &target, &cost, &nodes, &edges)
		if err != nil {
			panic(err)
		}

		g.tnrContractedNodes[[2]int64{source, target}] = nodes
		g.tnrContractedEdges[[2]int64{source, target}] = edges
		g.tnrContractedCosts[[2]int64{source, target}] = cost
	}

	rows.Close()
}

func loadTnrAccesses(g *Graphs, db *db.Postgres) {
	rows, err := db.Query(`SELECT junction, access_sources, access_sources_costs, access_targets, access_targets_costs FROM tnr_access`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var junction int64
		var accessSourceNodes, accessTargetNodes pq.Int64Array
		var accessSourceCosts, accessTargetCosts pq.Float64Array

		err := rows.Scan(&junction, &accessSourceNodes, &accessSourceCosts, &accessTargetNodes, &accessTargetCosts)
		if err != nil {
			panic(err)
		}

		g.tnrAccessSourceNodes[junction] = accessSourceNodes
		g.tnrAccessSourceCosts[junction] = accessSourceCosts
		g.tnrAccessTargetNodes[junction] = accessTargetNodes
		g.tnrAccessTargetCosts[junction] = accessTargetCosts
	}

	rows.Close()
}

func loadBaseGraph(g *ch.Graph, db *db.Postgres) {
	rows, err := db.Query(`SELECT id, source, target, cost, reverse_cost FROM base_graph`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, source, target int64
		var cost, reverseCost float64

		err := rows.Scan(&id, &source, &target, &cost, &reverseCost)
		if err != nil {
			panic(err)
		}

		g.CreateVertex(source)
		g.CreateVertex(target)

		if cost > 0 {
			g.AddEdge(source, target, cost)
		}

		if reverseCost > 0 {
			g.AddEdge(target, source, reverseCost)
		}
	}

	rows.Close()
}

func loadTnrGraph(g *ch.Graph, db *db.Postgres) {
	rows, err := db.Query(`SELECT id, source, target, cost, reverse_cost FROM tnr_graph`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, source, target int64
		var cost, reverseCost float64

		err := rows.Scan(&id, &source, &target, &cost, &reverseCost)
		if err != nil {
			panic(err)
		}

		g.CreateVertex(source)
		g.CreateVertex(target)

		if cost > 0 {
			g.AddEdge(source, target, cost)
		}

		if reverseCost > 0 {
			g.AddEdge(target, source, reverseCost)
		}
	}

	rows.Close()
}

func loadChGraph(g *ch.Graph, db *db.Postgres) {
	rows, err := db.Query(`SELECT source, target, cost FROM ch_graph`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var source, target int64
		var cost float64

		err := rows.Scan(&source, &target, &cost)
		if err != nil {
			panic(err)
		}

		g.CreateVertex(source)
		g.CreateVertex(target)
		g.AddEdge(source, target, cost)
	}

	rows.Close()

	rows, err = db.Query(`SELECT junction, orderpos, importance FROM ch_junctions`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var junction, orderpos, importance int64

		err := rows.Scan(&junction, &orderpos, &importance)
		if err != nil {
			panic(err)
		}

		vertexInternal, _ := g.FindVertex(junction)
		g.Vertices[vertexInternal].SetOrderPos(orderpos)
		g.Vertices[vertexInternal].SetImportance(int(importance))
	}

	rows.Close()

	rows, err = db.Query(`SELECT source, target, cost, via FROM ch_shortcuts`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var source, target, via int64
		var cost float64

		err := rows.Scan(&source, &target, &cost, &via)
		if err != nil {
			panic(err)
		}

		g.AddEdge(source, target, cost)
		g.AddShortcut(source, target, via, cost)
	}

	rows.Close()
}

func loadChTnrGraph(g *ch.Graph, db *db.Postgres) {
	rows, err := db.Query(`SELECT source, target, cost FROM ch_tnr_graph`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var source, target int64
		var cost float64

		err := rows.Scan(&source, &target, &cost)
		if err != nil {
			panic(err)
		}

		g.CreateVertex(source)
		g.CreateVertex(target)
		g.AddEdge(source, target, cost)
	}

	rows.Close()

	rows, err = db.Query(`SELECT junction, orderpos, importance FROM ch_tnr_junctions`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var junction, orderpos, importance int64

		err := rows.Scan(&junction, &orderpos, &importance)
		if err != nil {
			panic(err)
		}

		vertexInternal, _ := g.FindVertex(junction)
		g.Vertices[vertexInternal].SetOrderPos(orderpos)
		g.Vertices[vertexInternal].SetImportance(int(importance))
	}

	rows.Close()

	rows, err = db.Query(`SELECT source, target, cost, via FROM ch_tnr_shortcuts`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var source, target, via int64
		var cost float64

		err := rows.Scan(&source, &target, &cost, &via)
		if err != nil {
			panic(err)
		}

		g.AddEdge(source, target, cost)
		g.AddShortcut(source, target, via, cost)
	}

	rows.Close()
}
