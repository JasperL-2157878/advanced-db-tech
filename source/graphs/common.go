package graphs

import (
	"example.com/source/types"
)

type Graph interface {
	BdDijkstra(source int64, target int64) types.Path
}

type ContractionMap struct {
	data map[[2]int64][]int64
}

func (t *ContractionMap) Get(source int64, target int64) []int64 {
	if t.data == nil {
		t.data = make(map[[2]int64][]int64)
	}

	return t.data[[2]int64{source, target}]
}

func (t *ContractionMap) Set(source int64, target int64, id []int64) {
	if t.data == nil {
		t.data = make(map[[2]int64][]int64)
	}

	if len(id) > 0 {
		t.data[[2]int64{source, target}] = id
		t.data[[2]int64{target, source}] = id
	}
}
