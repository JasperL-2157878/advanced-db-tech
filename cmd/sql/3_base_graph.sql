SELECT
	nw.gid AS id,
	f_jnctid AS source, 
	t_jnctid AS target, 
	CASE
		WHEN COALESCE(nl.oneway, '') = 'FT' THEN minutes
		WHEN COALESCE(nl.oneway, '') = 'TF' THEN -1
		ELSE minutes
	END AS cost,
	CASE
		WHEN COALESCE(nl.oneway, '') = 'FT' THEN -1
		WHEN COALESCE(nl.oneway, '') = 'TF' THEN minutes
		ELSE minutes
	END AS reverse_cost,
	ST_X(ST_GeometryN(jcf.geom, 1)) AS x1,
	ST_Y(ST_GeometryN(jcf.geom, 1)) AS y1,
	ST_X(ST_GeometryN(jct.geom, 1)) AS x2,
	ST_Y(ST_GeometryN(jct.geom, 1)) AS y2
INTO base_graph
FROM nw 
	LEFT JOIN nl ON nw.id = nl.id
	LEFT JOIN jc jcf ON nw.f_jnctid = jcf.id
	LEFT JOIN jc jct ON nw.t_jnctid = jct.id;

CREATE INDEX fnw_idx_source ON fnw(source);
CREATE INDEX fnw_idx_target ON fnw(target);
