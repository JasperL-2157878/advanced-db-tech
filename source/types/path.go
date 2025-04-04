package types

import (
	"fmt"
)

type Path struct {
	Cost  float64
	Nodes []int64
	Edges []int64
}

func (p Path) ToTable() string {
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

func (p Path) NodesToArray() string {
	expr := "{"
	n := len(p.Nodes)

	for i, node := range p.Nodes {
		expr += fmt.Sprintf("%d", node)
		if i < n-1 {
			expr += ","
		}
	}

	return expr + "}"
}

func (p Path) EdgesToArray() string {
	expr := "{"
	n := len(p.Edges)

	for i, edge := range p.Edges {
		expr += fmt.Sprintf("%d", edge)
		if i < n-1 {
			expr += ","
		}
	}

	return expr + "}"
}
