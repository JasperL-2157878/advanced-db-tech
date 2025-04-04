package graphs

import (
	db "example.com/source/database"
	"example.com/source/types"
	"github.com/LdDl/ch"
)

type BaseGraph struct {
	super ch.Graph
	nodes ContractionMap
	edges ContractionMap
}

func LoadBaseGraph(db *db.Postgres) *BaseGraph {
	g := &BaseGraph{}

	rows, err := db.Query(`SELECT id, source, target, cost, reverse_cost FROM base_graph`)
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

	return g
}

func (g *BaseGraph) BdDijkstra(source int64, target int64) *types.Path {
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
