package graph

import (
	db "example.com/source/database"
	"github.com/lib/pq"
)

type Tnr struct {
	Shortcuts map[[2]int64]TnrShortcut
	Sources   map[int64][]TnrAccess
	Targets   map[int64][]TnrAccess
}

type TnrShortcut struct {
	Nodes []int64
	Edges []int64
	Cost  float64
}

type TnrAccess struct {
	Node int64
	Cost float64
}

func loadTnr(db *db.Postgres) *Tnr {
	tnr := &Tnr{
		Shortcuts: make(map[[2]int64]TnrShortcut),
		Sources:   make(map[int64][]TnrAccess),
		Targets:   make(map[int64][]TnrAccess),
	}

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

		tnr.Shortcuts[[2]int64{source, target}] = TnrShortcut{
			Nodes: nodes,
			Edges: edges,
			Cost:  cost,
		}
	}

	rows, err = db.Query(`SELECT junction, access_source_nodes, access_source_costs, access_target_nodes, access_target_costs FROM tnr_access`)
	if err != nil {
		panic(err)
	}

	tnr.Sources = make(map[int64][]TnrAccess)
	tnr.Targets = make(map[int64][]TnrAccess)

	for rows.Next() {
		var junction int64
		var sourceNodes, targetNodes pq.Int64Array
		var sourceCosts, targetCosts pq.Float64Array

		err := rows.Scan(&junction, &sourceNodes, &sourceCosts, &targetNodes, &targetCosts)
		if err != nil {
			panic(err)
		}

		for i := 0; i < 6; i++ {
			tnr.Sources[junction] = append(tnr.Sources[junction], TnrAccess{
				Node: sourceNodes[i],
				Cost: sourceCosts[i],
			})

			tnr.Targets[junction] = append(tnr.Targets[junction], TnrAccess{
				Node: targetNodes[i],
				Cost: targetCosts[i],
			})
		}
	}

	return tnr
}

/*func (g *Graphs) TnrDijkstra(source int64, target int64) *types.Path {
	start := time.Now()
	path := types.Path{}

	best := math.MaxFloat64
	tnrFrom := int64(-1)
	tnrTo := int64(-1)

	n := len(g.tnrSourceNodes[source])
	m := len(g.tnrTargetNodes[target])

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			dist := g.tnrSourceCosts[source][i] +
				g.tnrShortcutCosts[[2]int64{g.tnrSourceNodes[source][i], g.tnrTargetNodes[target][j]}] +
				g.tnrTargetCosts[target][j]

			if dist < best {
				best = dist
				tnrFrom = g.tnrSourceNodes[source][i]
				tnrTo = g.tnrTargetNodes[target][j]
			}
		}
	}

	g.tnrAppendPathVanillaDijkstra(&path, source, tnrFrom)
	path.Nodes = append(path.Nodes, g.tnrShortcutNodes[[2]int64{tnrFrom, tnrTo}]...)
	path.Edges = append(path.Edges, g.tnrShortcutEdges[[2]int64{tnrFrom, tnrTo}]...)
	path.Cost += g.tnrShortcutCosts[[2]int64{tnrFrom, tnrTo}]
	g.tnrAppendPathVanillaDijkstra(&path, tnrTo, target)
	path.Edges = append(path.Edges, -1)
	path.QueryTime = time.Since(start)

	return &path
}

func (g *Graphs) tnrAppendPathVanillaDijkstra(path *types.Path, source int64, target int64) {
	cost, nodes := g.Base.VanillaShortestPath(source, target)
	n := len(nodes)

	fmt.Println(source, target, cost, nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Edges = append(path.Edges, g.edges[[2]int64{source, target}])
	}

	path.Cost += cost
	path.Nodes = append(path.Nodes, nodes...)
}
*/
