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
