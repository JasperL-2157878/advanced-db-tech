package graphs

import (
	"log"

	db "example.com/source/database"
	"example.com/source/types"
)

type Graph interface {
	BdDijkstra(source int64, target int64) types.Path
}

type Graphs struct {
	Base  *BaseGraph
	Tnr   *TnrGraph
	Ch    *ChGraph
	ChTnr *ChTnrGraph
}

func LoadGraphs(db *db.Postgres) *Graphs {
	log.Println("[INFO] Loading graphs into memory ...")
	return &Graphs{
		Base:  LoadBaseGraph(db),
		Tnr:   LoadTnrGraph(db),
		Ch:    LoadChGraph(db),
		ChTnr: LoadChTnrGraph(db),
	}
}
