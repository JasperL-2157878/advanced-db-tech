package db

func (pg *PostgresConnection) Geocode(street string, number int, postal string, city string) []byte {
	query := pg.conn.QueryRow(`
		SELECT json_agg(
		  json_build_object(
			'id',
			id,
			'fullname',
			fullname,
			'l_axon',
			l_axon,
			'l_pc',
			l_pc,
			'r_axon',
			r_axon,
			'r_pc',
			r_pc,
			'l_f_add',
			l_f_add,
			'l_t_add',
			l_t_add,
			'r_f_add',
			r_f_add,
			'r_t_add',
			r_t_add,
			'f_jnctid',
			f_jnctid,
			't_jnctid',
		    t_jnctid
 		  )
		)
		FROM (
		  SELECT DISTINCT ON (gc.fullname, gc.l_axon, gc.l_f_add, gc.l_t_add, gc.r_f_add, gc.r_t_add)
		    gc.id,
		    gc.fullname,
		    gc.l_axon,
		    gc.l_pc,
		    gc.r_axon,
		    gc.r_pc,
			gc.l_f_add,
			gc.l_t_add,
			gc.r_f_add,
			gc.r_t_add,
		    nw.f_jnctid,
		    nw.t_jnctid
		  FROM gc JOIN nw ON gc.id = nw.id
		  WHERE 
		    UPPER(gc.fullname) LIKE UPPER($1) AND (
			  $2 = 0 OR $2 BETWEEN LEAST(
			    gc.l_f_add, gc.l_t_add, gc.r_f_add, gc.r_t_add
			  ) AND GREATEST(
			    gc.l_f_add, gc.l_t_add, gc.r_f_add, gc.r_t_add
			  )
		    )
		    AND 
		    ($3 = '' OR ($3 = gc.l_pc OR $3 = gc.r_pc))
		    AND
		    ($4 = '' OR (UPPER($4) = UPPER(gc.l_axon) OR UPPER($4) = UPPER(gc.r_axon)))
		  ORDER BY gc.fullname, gc.l_axon ASC
		  LIMIT 10
		)
	`, street+"%", number, postal, city)

	var json []byte
	err := query.Scan(&json)
	if err != nil {
		panic(err)
	}

	return json
}
