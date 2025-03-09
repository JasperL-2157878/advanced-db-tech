package db

func (pg *PostgresConnection) Example() []byte {
	query := pg.conn.QueryRow(`
		SELECT json_build_object(
			'type', 'FeatureCollection',
			'features', json_agg(
				json_build_object(
					'type', 'Feature',
					'geometry', ST_AsGeoJSON(geom)::json,
					'properties', json_build_object('gid', gid)
				)
			)
		) AS geojson
		FROM nw
		WHERE gid IN (
			SELECT edge FROM pgr_dijkstra(
				'SELECT gid AS id, f_jnctid AS source, t_jnctid AS target, meters AS cost, meters AS reverse_cost FROM nw',
				10560298937508,
				10560298942473,
				false
			)
		)
	`)

	var json []byte
	query.Scan(&json)

	return json
}
