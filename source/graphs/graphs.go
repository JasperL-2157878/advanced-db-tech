package graphs

import (
	"log"

	db "example.com/source/database"
	"github.com/LdDl/ch"
)

type Graphs struct {
	// Different graphs
	Base  *ch.Graph
	Tnr   *ch.Graph
	Ch    *ch.Graph
	ChTnr *ch.Graph

	// Shared memory
	edges map[[2]int64]int64

	tnrContractedNodes map[[2]int64][]int64
	tnrContractedEdges map[[2]int64][]int64
	tnrContractedCosts map[[2]int64]float64

	tnrAccessSourceNodes map[int64][]int64
	tnrAccessSourceCosts map[int64][]float64
	tnrAccessTargetNodes map[int64][]int64
	tnrAccessTargetCosts map[int64][]float64
}

func LoadGraphs(db *db.Postgres) *Graphs {
	log.Println("[INFO] Loading graphs into memory ...")

	g := &Graphs{
		Base:  &ch.Graph{},
		Tnr:   &ch.Graph{},
		Ch:    &ch.Graph{},
		ChTnr: &ch.Graph{},

		edges: make(map[[2]int64]int64),

		tnrContractedNodes: make(map[[2]int64][]int64),
		tnrContractedEdges: make(map[[2]int64][]int64),
		tnrContractedCosts: make(map[[2]int64]float64),

		tnrAccessSourceNodes: make(map[int64][]int64),
		tnrAccessSourceCosts: make(map[int64][]float64),
		tnrAccessTargetNodes: make(map[int64][]int64),
		tnrAccessTargetCosts: make(map[int64][]float64),
	}

	loadEdges(g, db)
	//loadTnrShortcuts(g, db)
	//loadTnrAccesses(g, db)

	//loadBaseGraph(g.Base, db)
	//loadTnrGraph(g.Tnr, db)
	loadChGraph(g.Ch, db)
	//loadChTnrGraph(g.ChTnr, db)

	return g
}
