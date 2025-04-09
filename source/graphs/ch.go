package graphs

import (
	"time"

	"example.com/source/types"
)

func (g *Graphs) ChDijkstra(source int64, target int64) *types.Path {
	start := time.Now()
	path := types.Path{}

	cost, nodes := g.Ch.ShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Edges = append(path.Edges, g.edges[[2]int64{source, target}])
	}

	path.Cost = cost
	path.Nodes = nodes
	path.QueryTime = time.Since(start)

	return &path
}
