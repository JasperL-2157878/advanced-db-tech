package graphs

import (
	"math"
	"time"

	db "example.com/source/database"
	"example.com/source/types"
	"github.com/LdDl/ch"
	"github.com/lib/pq"
)

type TnrGraph struct {
	super           ch.Graph
	contractedNodes map[[2]int64][]int64
	contractedEdges map[[2]int64][]int64
	accessNodes     map[int64][]int64
	accessCosts     map[int64][]float64
	transitCosts    map[[2]int64]float64
}

func LoadTnrGraph(db *db.Postgres) *TnrGraph {
	g := &TnrGraph{}
	g.contractedNodes = make(map[[2]int64][]int64)
	g.contractedEdges = make(map[[2]int64][]int64)
	g.accessNodes = make(map[int64][]int64)
	g.accessCosts = make(map[int64][]float64)
	g.transitCosts = make(map[[2]int64]float64)

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

	rows, err = db.Query(`SELECT source, target, cost, nodes, edges FROM tnr_shortcuts`)
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

		g.contractedNodes[[2]int64{source, target}] = nodes
		g.contractedNodes[[2]int64{target, source}] = nodes
		g.contractedEdges[[2]int64{source, target}] = edges
		g.contractedEdges[[2]int64{target, source}] = edges
		g.transitCosts[[2]int64{source, target}] = cost
		g.transitCosts[[2]int64{target, source}] = cost
	}

	rows.Close()

	rows, err = db.Query(`SELECT junction, access, cost FROM tnr_access`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var junction int64
		var access pq.Int64Array
		var costs pq.Float64Array

		err := rows.Scan(&junction, &access, &costs)
		if err != nil {
			panic(err)
		}

		g.accessNodes[junction] = access
		g.accessCosts[junction] = costs
	}

	rows.Close()

	return g
}

func (g *TnrGraph) BdDijkstra(source int64, target int64) *types.Path {
	start := time.Now()
	path := types.Path{}

	_, sourceExists := g.accessNodes[source]
	_, targetExists := g.accessNodes[target]
	if !sourceExists || !targetExists {
		g.appendPath(&path, source, target)
		return &path
	}

	best := math.MaxFloat64
	sourceAccess := int64(-1)
	targetAccess := int64(-1)

	n := len(g.accessNodes[source])
	m := len(g.accessNodes[target])

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			dist := g.accessCosts[source][i] +
				g.transitCosts[[2]int64{source, target}] +
				g.accessCosts[target][j]

			if dist < best {
				best = dist
				sourceAccess = g.accessNodes[source][i]
				targetAccess = g.accessNodes[target][j]
			}
		}
	}

	g.appendPath(&path, source, sourceAccess)
	path.Nodes = append(path.Nodes, g.contractedNodes[[2]int64{sourceAccess, targetAccess}]...)
	path.Edges = append(path.Edges, g.contractedEdges[[2]int64{sourceAccess, targetAccess}]...)
	path.Cost += g.transitCosts[[2]int64{sourceAccess, targetAccess}]
	g.appendPath(&path, targetAccess, target)

	path.QueryTime = time.Since(start)

	return &path
}

func (g *TnrGraph) appendPath(path *types.Path, source int64, target int64) {
	cost, nodes := g.super.VanillaShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Nodes = append(path.Nodes, g.contractedNodes[[2]int64{source, target}]...)
		path.Edges = append(path.Edges, g.contractedEdges[[2]int64{source, target}]...)
	}

	path.Cost += cost
}
