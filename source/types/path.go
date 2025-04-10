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
	startTime time.Time
}

func NewPath() Path {
	return Path{
		Cost:      0,
		Nodes:     make([]int64, 0),
		Edges:     make([]int64, 0),
		QueryTime: time.Duration(0),
		startTime: time.Now(),
	}
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

func (p *Path) Append(q *Path) {
	p.Cost += q.Cost
	p.Nodes = append(p.Nodes, q.Nodes...)
	p.Edges = append(p.Edges, q.Edges...)
}

func (p *Path) End() {
	p.QueryTime = time.Since(p.startTime)
	p.Edges = append(p.Edges, -1)
}
