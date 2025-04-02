package preprocessors

import (
	db "example.com/source/database"
	"github.com/LdDl/ch"
	_ "github.com/lib/pq"
)

func GenerateContractionHierarchies(db *db.Postgres) {
	g := load(db)

	g.PrepareContractionHierarchies()
	g.ExportToFile("ch")
}

func load(db *db.Postgres) *ch.Graph {
	g := &ch.Graph{}

	rows, err := db.Query(`
		SELECT
			f_jnctid AS source, 
			t_jnctid AS target,
			COALESCE(nl.oneway, '') AS oneway,
		CASE
			WHEN COALESCE(nl.oneway, '') = 'FT' THEN minutes
			WHEN COALESCE(nl.oneway, '') = 'TF' THEN -1
			ELSE minutes
		END AS cost,
		CASE
			WHEN COALESCE(nl.oneway, '') = 'FT' THEN -1
			WHEN COALESCE(nl.oneway, '') = 'TF' THEN minutes
			ELSE minutes
		END AS reverse_cost
		FROM nw LEFT JOIN nl ON nw.id = nl.id
	`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var source, target int64
		var oneway string
		var cost, reverseCost float64

		err := rows.Scan(&source, &target, &oneway, &cost, &reverseCost)
		if err != nil {
			panic(err)
		}

		err = g.CreateVertex(source)
		if err != nil {
			panic(err)
		}
		err = g.CreateVertex(target)
		if err != nil {
			panic(err)
		}

		if (oneway == "FT" || oneway == "N" || oneway == "") && cost > 0 {
			err = g.AddEdge(source, target, cost)
			if err != nil {
				panic(err)
			}
		}

		if (oneway == "TF" || oneway == "N" || oneway == "") && reverseCost > 0 {
			err = g.AddEdge(target, source, reverseCost)
			if err != nil {
				panic(err)
			}
		}
	}

	return g
}
