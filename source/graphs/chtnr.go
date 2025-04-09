package graphs

import (
	"math"
	"time"

	"example.com/source/types"
	"github.com/LdDl/ch"
)

func (g *Graphs) ChTnrDijkstra(source int64, target int64) *types.Path {
	start := time.Now()
	path := types.Path{}

	best := math.MaxFloat64
	accessSourceNode := int64(-1)
	accessTargetNode := int64(-1)

	n := len(g.tnrAccessSourceNodes[source])
	m := len(g.tnrAccessTargetNodes[target])

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			dist := g.tnrAccessSourceCosts[source][i] +
				g.tnrContractedCosts[[2]int64{g.tnrAccessSourceNodes[source][i], g.tnrAccessTargetNodes[target][j]}] +
				g.tnrAccessTargetCosts[target][j]

			if dist < best {
				best = dist
				accessSourceNode = g.tnrAccessSourceNodes[source][i]
				accessTargetNode = g.tnrAccessTargetNodes[target][j]
			}
		}
	}

	g.tnrAppendExtendedDijkstra(&path, g.Tnr, source, accessSourceNode)
	path.Nodes = append(path.Nodes, g.tnrContractedNodes[[2]int64{accessSourceNode, accessTargetNode}]...)
	path.Edges = append(path.Edges, g.tnrContractedEdges[[2]int64{accessSourceNode, accessTargetNode}]...)
	path.Cost += g.tnrContractedCosts[[2]int64{accessSourceNode, accessTargetNode}]
	g.tnrAppendExtendedDijkstra(&path, g.Tnr, accessTargetNode, target)

	path.QueryTime = time.Since(start)

	return &path
}

func (g *Graphs) tnrAppendExtendedDijkstra(path *types.Path, graph *ch.Graph, source int64, target int64) {
	cost, nodes := graph.ShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Nodes = append(path.Nodes, g.tnrContractedNodes[[2]int64{source, target}]...)
		path.Edges = append(path.Edges, g.tnrContractedNodes[[2]int64{source, target}]...)
	}

	path.Cost += cost
}
