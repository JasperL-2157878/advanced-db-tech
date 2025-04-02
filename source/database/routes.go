package db

import (
	"encoding/json"
	"fmt"
)

func (db *Postgres) Route(path Path) GeoJSON {
	query := db.QueryRow(fmt.Sprintf(`
		WITH route AS (
			SELECT 
				path.seq,
				path.agg_cost,
				nw.gid,
				nw.name,
				nw.meters,
				nw.routenum,
				nl.minutes,
				nl.fow,
				CASE
					WHEN path.node = nw.f_jnctid THEN nw.geom
					ELSE ST_Reverse(nw.geom)
				END as geom
			FROM nw JOIN %s ON nw.gid = path.edge
			LEFT JOIN nl ON nw.id = nl.id
			ORDER BY path.seq
		),
		include_previous_route AS (
			SELECT 
				*,
				ST_LineMerge(geom) as line_geom,
				LAG(ST_LineMerge(geom)) OVER (ORDER BY seq) AS prev_line_geom
			FROM route
		),
		angles AS (
			SELECT 
				*,
				CAST(degrees(ST_Azimuth(
					ST_StartPoint(line_geom),
					CASE
						WHEN ST_NPoints(line_geom) = 2 THEN ST_PointN(line_geom, 2)
						ELSE ST_PointN(line_geom, 3)
					END
				)) AS numeric) AS current_angle,
			CAST(degrees(ST_Azimuth(
				CASE
					WHEN ST_NPoints(prev_line_geom) = 2 THEN ST_PointN(prev_line_geom, ST_NPoints(prev_line_geom) - 1)
					ELSE ST_PointN(prev_line_geom, ST_NPoints(prev_line_geom) - 2)
				END,
				ST_EndPoint(prev_line_geom)
				)) AS numeric) AS prev_angle
			FROM include_previous_route
		),
		diff_angles AS (
			SELECT
				*,
				CASE 
					WHEN prev_angle IS NULL THEN NULL
					ELSE (current_angle - prev_angle + 180 + 360) %% 360 - 180
				END AS angle_diff
			FROM angles
		)
		SELECT json_build_object(
			'type', 'FeatureCollection',
			'total_cost', (SELECT agg_cost FROM route ORDER BY seq DESC LIMIT 1),
			'features', json_agg(
				json_build_object(
					'type', 'Feature',
					'geometry', ST_AsGeoJSON(geom)::json,
					'properties', json_build_object(
						'gid', gid,
						'street_name', name,
						'route_num', routenum,
						'fow', fow,
						'angle_diff', angle_diff,
						'distance', meters,
						'duration', minutes
					)
				) ORDER BY seq
			)
		) AS geojson
		FROM diff_angles;
	`, path.ToTable()))

	var rawJson []byte
	err := query.Scan(&rawJson)
	if err != nil {
		panic(err)
	}

	var jsonValue GeoJSON
	err = json.Unmarshal(rawJson, &jsonValue)
	if err != nil {
		panic(err)
	}

	return jsonValue
}

func (db *Postgres) path(alg string, from int64, to int64) Path {
	rows, err := db.Query(fmt.Sprintf(`
		SELECT seq, node, edge, cost FROM pgr_%s(
			'SELECT * FROM fnw',
			CAST($1 AS BIGINT),
			CAST($2 AS BIGINT),
			true
		)
	`, alg), from, to)

	if err != nil {
		return Path{}
	}

	var path Path
	for rows.Next() {
		var segment PathSegment
		err := rows.Scan(&segment.Seq, &segment.Node, &segment.Edge, &segment.AggCost)
		if err != nil {
			panic(err)
		}

		path.Sequences = append(path.Sequences, segment)
	}

	return path
}

func (db *Postgres) Dijkstra(from int64, to int64) Path {
	return db.path("dijkstra", from, to)
}

func (db *Postgres) Astar(from int64, to int64) Path {
	return db.path("aStar", from, to)
}

func (db *Postgres) BdDijkstra(from int64, to int64) Path {
	return db.path("bdDijkstra", from, to)
}

func (db *Postgres) BdAstar(from int64, to int64) Path {
	return db.path("bdAstar", from, to)
}
