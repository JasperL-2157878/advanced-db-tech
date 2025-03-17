package db

import "fmt"

func (pg *PostgresConnection) Example() []byte {
	query := pg.conn.QueryRow(`
		WITH route AS (
		  SELECT 
		      pgr.seq,
		      
		      nw.gid,
		      nw.name,
		      nw.meters,
		      nl.fow,
		      
		      CASE
		          WHEN pgr.node = nw.f_jnctid THEN nw.geom
		          ELSE ST_Reverse(nw.geom)
		      END as geom
		  FROM nw
		  JOIN (
		    SELECT * FROM pgr_dijkstra(
		      'SELECT 
				nw.gid AS id, 
				f_jnctid AS source, 
				t_jnctid AS target, 
				CASE
				  WHEN COALESCE(nl.oneway, '''') = ''FT'' THEN minutes
				  WHEN COALESCE(nl.oneway, '''') = ''TF'' THEN -1
				  ELSE minutes
				END AS cost,
				CASE
				  WHEN COALESCE(nl.oneway, '''') = ''FT'' THEN -1
				  WHEN COALESCE(nl.oneway, '''') = ''TF'' THEN minutes
				  ELSE minutes
				END AS reverse_cost
				FROM nw LEFT JOIN nl ON nw.id = nl.id',
		      10560298937508,
		      10560298942473,
		      true
		    )
		  ) AS pgr ON nw.gid = pgr.edge
		  LEFT JOIN nl ON nw.id = nl.id
		  ORDER BY pgr.seq
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
		      ELSE (current_angle - prev_angle + 180 + 360) % 360 - 180
		    END AS angle_diff
		  FROM angles
		)
		SELECT json_build_object(
		  'type', 'FeatureCollection',
		  'features', json_agg(
		    json_build_object(
		      'type', 'Feature',
		      'geometry', ST_AsGeoJSON(geom)::json,
		      'properties', json_build_object(
		        'gid', gid,
		        'street_name', name,
		        'distance', meters,
		        'fow', fow,
		        'angle_diff', angle_diff
		      )
		    ) ORDER BY seq
		  )
		) AS geojson
		FROM diff_angles;
	`)

	var json []byte
	err := query.Scan(&json)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return json
}
