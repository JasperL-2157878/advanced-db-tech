package graph

import (
	"log"
	"math"

	db "example.com/source/database"
	"example.com/source/types"
	"github.com/LdDl/ch"
)

var edges map[[2]int64]int64
var tnr *Tnr

type Graph struct {
	ch.Graph
}

func LoadGraph(db *db.Postgres) *Graph {
	log.Println("[INFO] Loading graph into memory (~7.5GB) ...")

	edges = make(map[[2]int64]int64)
	tnr = loadTnr(db)

	g := Graph{}

	rows, err := db.Query(`SELECT id, source, target FROM base_graph`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, source, target int64

		err := rows.Scan(&id, &source, &target)
		if err != nil {
			panic(err)
		}

		edges[[2]int64{source, target}] = id
		edges[[2]int64{target, source}] = id
	}

	rows, err = db.Query(`SELECT source, target, cost FROM ch_graph`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var source, target int64
		var cost float64

		err := rows.Scan(&source, &target, &cost)
		if err != nil {
			panic(err)
		}

		g.CreateVertex(source)
		g.CreateVertex(target)
		g.AddEdge(source, target, cost)
	}

	rows, err = db.Query(`SELECT junction, order_pos, importance FROM ch_junctions`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var junction, orderPos, importance int64

		err := rows.Scan(&junction, &orderPos, &importance)
		if err != nil {
			panic(err)
		}

		vertexInternal, _ := g.FindVertex(junction)
		g.Vertices[vertexInternal].SetOrderPos(orderPos)
		g.Vertices[vertexInternal].SetImportance(int(importance))
	}

	rows, err = db.Query(`SELECT source, target, cost, via FROM ch_shortcuts`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var source, target, via int64
		var cost float64

		err := rows.Scan(&source, &target, &cost, &via)
		if err != nil {
			panic(err)
		}

		g.AddEdge(source, target, cost)
		g.AddShortcut(source, target, via, cost)
	}

	return &g
}

func (g *Graph) Base(source, target int64) *types.Path {
	path := types.NewPath()

	cost, nodes := g.VanillaShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Edges = append(path.Edges, edges[[2]int64{source, target}])
	}

	path.Cost = cost
	path.Nodes = nodes
	path.End()

	return &path
}

func (g *Graph) Ch(source, target int64) *types.Path {
	path := types.NewPath()

	cost, nodes := g.ShortestPath(source, target)
	n := len(nodes)

	for i := 1; i < n; i++ {
		source := nodes[i-1]
		target := nodes[i]

		path.Edges = append(path.Edges, edges[[2]int64{source, target}])
	}

	path.Cost = cost
	path.Nodes = nodes
	path.End()

	return &path
}

func (g *Graph) BaseTnr(source, target int64) *types.Path {
	path := types.NewPath()

	best := math.MaxFloat64
	tnrFrom := int64(-1)
	tnrTo := int64(-1)

	n := len(tnr.Sources[source])
	m := len(tnr.Targets[target])

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			shortcut, exists := tnr.Shortcuts[[2]int64{tnr.Sources[source][i].Node, tnr.Targets[target][j].Node}]
			if !exists {
				continue
			}

			dist := tnr.Sources[source][i].Cost + shortcut.Cost + tnr.Targets[target][j].Cost

			if dist < best {
				best = dist
				tnrFrom = tnr.Sources[source][i].Node
				tnrTo = tnr.Targets[target][j].Node
			}
		}
	}

	if tnrTo < 0 || tnrFrom < 0 {
		return g.Base(source, target)
	}

	path.Append(g.Base(source, tnrFrom))
	path.Append(&types.Path{
		Nodes: tnr.Shortcuts[[2]int64{tnrFrom, tnrTo}].Nodes,
		Edges: tnr.Shortcuts[[2]int64{tnrFrom, tnrTo}].Edges,
		Cost:  tnr.Shortcuts[[2]int64{tnrFrom, tnrTo}].Cost,
	})
	path.Append(g.Base(tnrTo, target))
	path.End()

	return &path
}

func (g *Graph) ChTnr(source, target int64) *types.Path {
	path := types.NewPath()

	best := math.MaxFloat64
	tnrFrom := int64(-1)
	tnrTo := int64(-1)

	n := len(tnr.Sources[source])
	m := len(tnr.Targets[target])

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			shortcut, exists := tnr.Shortcuts[[2]int64{tnr.Sources[source][i].Node, tnr.Targets[target][j].Node}]
			if !exists {
				continue
			}

			dist := tnr.Sources[source][i].Cost + shortcut.Cost + tnr.Targets[target][j].Cost

			if dist < best {
				best = dist
				tnrFrom = tnr.Sources[source][i].Node
				tnrTo = tnr.Targets[target][j].Node
			}
		}
	}

	if tnrTo < 0 || tnrFrom < 0 {
		return g.Ch(source, target)
	}

	path.Append(g.Ch(source, tnrFrom))
	path.Append(&types.Path{
		Nodes: tnr.Shortcuts[[2]int64{tnrFrom, tnrTo}].Nodes,
		Edges: tnr.Shortcuts[[2]int64{tnrFrom, tnrTo}].Edges,
		Cost:  tnr.Shortcuts[[2]int64{tnrFrom, tnrTo}].Cost,
	})
	path.Append(g.Ch(tnrTo, target))
	path.End()

	return &path
}
