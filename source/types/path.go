package types

import (
	"fmt"
	"time"
)

type Path struct {
	Cost      float64
	Nodes     []int64
	Edges     []int64
	QueryTime time.Duration
}

func (p *Path) ToTable() string {
	expr := "(VALUES "
	n := len(p.Nodes)

	for i := 0; i < n; i++ {
		expr += fmt.Sprintf("(%d,%d,%d)", i+1, p.Nodes[i], p.Edges[i])
		if i < n-1 {
			expr += ","
		}
	}

	return expr + ") AS path(seq,node,edge)"
}
