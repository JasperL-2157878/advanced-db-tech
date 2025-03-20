package db

import (
	"regexp"
	"strconv"
)

func (pg *PostgresConnection) Geocode(address string) []byte {
	street, number, postal, city := pg.parseAddress(address)
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
		  SELECT DISTINCT ON (gc.fullname, gc.l_axon)
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
			  $2 = 0 OR (gc.l_f_add != -1 AND gc.r_f_add != -1 AND $2 BETWEEN LEAST(
			    gc.l_f_add, gc.l_t_add, gc.r_f_add, gc.r_t_add
			  ) AND GREATEST(
			    gc.l_f_add, gc.l_t_add, gc.r_f_add, gc.r_t_add
			  ))
		    )
		    AND 
		    ($3 = '' OR (gc.l_pc LIKE $3 OR gc.r_pc LIKE $3))
		    AND
		    ($4 = '' OR (UPPER(gc.l_axon) LIKE UPPER($4) OR UPPER(gc.r_axon) LIKE UPPER($4)))
		  ORDER BY gc.fullname, gc.l_axon, gc.id ASC
		  LIMIT 10
		)
	`, street+"%", number, postal+"%", city+"%")

	var json []byte
	err := query.Scan(&json)
	if err != nil {
		panic(err)
	}

	if len(json) == 0 {
		return []byte("[]")
	}

	return json
}

func (pg *PostgresConnection) parseAddress(address string) (street string, number int, postal string, city string) {
	re := regexp.MustCompile(`^(?P<street>[^0-9,]+)\s*(?P<number>\d+)?[,\s]*(?P<postal>\d{4})?\s*(?P<city>\D+)?$`)

	matches := re.FindStringSubmatch(address)
	street = matches[re.SubexpIndex("street")]
	number, _ = strconv.Atoi(matches[re.SubexpIndex("number")])
	postal = matches[re.SubexpIndex("postal")]
	city = matches[re.SubexpIndex("city")]

	println("street:", "\""+street+"\"")
	println("number:", number)
	println("postal:", "\""+postal+"\"")
	println("city:", "\""+city+"\"")

	return
}
