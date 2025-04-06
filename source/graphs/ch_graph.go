package graphs

import (
	"time"

	db "example.com/source/database"
	"example.com/source/types"
	"github.com/LdDl/ch"
)

type ChGraph struct {
	super           ch.Graph
	contractedNodes map[[2]int64][]int64
	contractedEdges map[[2]int64][]int64
}

func LoadChGraph(db *db.Postgres) *ChGraph {
	g := &ChGraph{}
	g.contractedNodes = make(map[[2]int64][]int64)
	g.contractedEdges = make(map[[2]int64][]int64)

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

		g.contractedNodes[[2]int64{source, target}] = []int64{source}
		g.contractedNodes[[2]int64{target, source}] = []int64{source}
		g.contractedEdges[[2]int64{source, target}] = []int64{id}
		g.contractedEdges[[2]int64{target, source}] = []int64{id}
	}

	rows.Close()

	rows, err = db.Query(`SELECT source, target, cost FROM ch_graph`)
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

		g.super.CreateVertex(source)
		g.super.CreateVertex(target)
		g.super.AddEdge(source, target, cost)
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

		vertexInternal, _ := g.super.FindVertex(junction)
		g.super.Vertices[vertexInternal].SetOrderPos(orderpos)
		g.super.Vertices[vertexInternal].SetImportance(int(importance))
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

		g.super.AddEdge(source, target, cost)
		g.super.AddShortcut(source, target, via, cost)
	}

	rows.Close()

	return g
}

func (g *ChGraph) BdDijkstra(source int64, target int64) *types.Path {
	start := time.Now()
	path := types.Path{}

	cost, nodes := g.super.ShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Nodes = append(path.Nodes, g.contractedNodes[[2]int64{source, target}]...)
		path.Edges = append(path.Edges, g.contractedEdges[[2]int64{source, target}]...)
	}

	path.Cost = cost
	path.QueryTime = time.Since(start)

	return &path
}
