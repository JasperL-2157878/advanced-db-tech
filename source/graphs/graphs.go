package graphs

import (
	db "example.com/source/database"
)

type Graphs struct {
	Base *BaseGraph
	Tnr  *TnrGraph
}

func LoadGraphs(db *db.Postgres) *Graphs {
	return &Graphs{
		Base: LoadBaseGraph(db),
		Tnr:  LoadTnrGraph(db),
	}
}
