package db

import (
	"example.com/source/types"
	"fmt"
)

func (db *Postgres) Places(input string) types.JSON {
	query := db.QueryRow(`
        SELECT
            json_agg(
                json_build_object(
                    'id', pi.cltrpelid,
                    'fullname', pi.name,
                	'l_axon', gc.l_axon,
                	'l_pc', gc.l_pc,
                	'r_axon', gc.r_axon,
                	'r_pc', gc.r_pc,
                	'f_jnctid', nw.f_jnctid,
                	't_jnctid', nw.t_jnctid
                )
            )
        FROM pi JOIN gc ON pi.cltrpelid = gc.id JOIN nw ON gc.id = nw.id
        WHERE pi.name ILIKE CONCAT('%', $1::TEXT, '%') LIMIT 10
	`, input)

	var json []byte
	err := query.Scan(&json)
	if err != nil {
		panic(err)
	}

	// dump the json to the console
	fmt.Println(string(json))

	if len(json) == 0 {
		return types.JSON("[]")
	}

	return json
}
