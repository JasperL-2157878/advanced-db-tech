package graphs

import (
	db "example.com/source/database"
	"example.com/source/types"
	"github.com/LdDl/ch"
)

type BaseGraph struct {
	super           ch.Graph
	contractedNodes map[[2]int64][]int64
	contractedEdges map[[2]int64][]int64
}

func LoadBaseGraph(db *db.Postgres) *BaseGraph {
	g := &BaseGraph{}
	g.contractedNodes = make(map[[2]int64][]int64)
	g.contractedEdges = make(map[[2]int64][]int64)

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

		g.super.CreateVertex(source)
		g.super.CreateVertex(target)

		if cost > 0 {
			g.super.AddEdge(source, target, cost)
		}

		if reverseCost > 0 {
			g.super.AddEdge(target, source, reverseCost)
		}

		g.contractedNodes[[2]int64{source, target}] = []int64{source}
		g.contractedNodes[[2]int64{target, source}] = []int64{source}
		g.contractedEdges[[2]int64{source, target}] = []int64{id}
		g.contractedEdges[[2]int64{target, source}] = []int64{id}
	}

	rows.Close()

	return g
}

func (g *BaseGraph) BdDijkstra(source int64, target int64) *types.Path {
	path := types.Path{}

	cost, nodes := g.super.VanillaShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Nodes = append(path.Nodes, g.contractedNodes[[2]int64{source, target}]...)
		path.Edges = append(path.Edges, g.contractedEdges[[2]int64{source, target}]...)
	}

	path.Cost = cost

	return &path
}
