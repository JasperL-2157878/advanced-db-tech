package graphs

import (
	db "example.com/source/database"
	"example.com/source/types"
	"github.com/LdDl/ch"
	"github.com/lib/pq"
)

type TnrGraph struct {
	super ch.Graph
	nodes ContractionMap
	edges ContractionMap
}

func LoadTnrGraph(db *db.Postgres) *TnrGraph {
	g := &TnrGraph{}

	rows, err := db.Query(`SELECT id, source, target, cost, reverse_cost FROM tnr_graph`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, source, target int64
		var cost, reverseCost float64

		err := rows.Scan(&id, &source, &target, &cost, &reverseCost)
		if err != nil {
			panic(err)
		}

		g.super.CreateVertex(source)
		g.super.CreateVertex(target)

		if cost > 0 {
			g.super.AddEdge(source, target, cost)
		}

		if reverseCost > 0 {
			g.super.AddEdge(target, source, reverseCost)
		}

		g.nodes.Set(source, target, []int64{source})
		g.edges.Set(source, target, []int64{id})
	}

	rows, err = db.Query(`SELECT source, target, nodes, edges FROM tnr_shortcuts`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var source, target int64
		var nodes, edges pq.Int64Array

		err := rows.Scan(&source, &target, &nodes, &edges)
		if err != nil {
			panic(err)
		}

		g.nodes.Set(source, target, nodes)
		g.edges.Set(source, target, edges)
	}

	return g
}

func (g *TnrGraph) BdDijkstra(source int64, target int64) *types.Path {
	path := types.Path{}

	cost, nodes := g.super.VanillaShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Nodes = append(path.Nodes, g.nodes.Get(source, target)...)
		path.Edges = append(path.Edges, g.edges.Get(source, target)...)
	}

	path.Cost = cost

	return &path
}
