CREATE TABLE tnr_shortcuts (
  id BIGINT,
  source BIGINT,
  target BIGINT,
  cost DOUBLE PRECISION,
  reverse_cost DOUBLE PRECISION,
  via BIGINT[]
);

COPY tnr_shortcuts(id, source, target, cost, reverse_cost, via)
FROM '/docker-entrypoint-initdb.d/preprocessings/tnr_shortcuts.csv'
DELIMITER ';'
CSV HEADER;

CREATE INDEX tnr_shortcuts_idx_source ON tnr_shortcuts(source);
CREATE INDEX tnr_shortcuts_idx_target ON tnr_shortcuts(target);

SELECT * INTO tnr_graph
FROM base_graph
WHERE NOT EXISTS (
    SELECT 1
    FROM tnr_shortcuts
    WHERE
      base_graph.id = ANY(tnr_shortcuts.edges)
);

CREATE INDEX tnr_graph_idx_source ON tnr_graph(source);
CREATE INDEX tnr_graph_idx_target ON tnr_graph(target);

INSERT INTO tnr_graph
SELECT
    tnr_shortcuts.id, source, target, cost, reverse_cost,
    ST_X(ST_GeometryN(jcf.geom, 1)) AS x1,
    ST_Y(ST_GeometryN(jcf.geom, 1)) AS y1,
    ST_X(ST_GeometryN(jct.geom, 1)) AS x2,
    ST_Y(ST_GeometryN(jct.geom, 1)) AS y2
FROM tnr_shortcuts
	LEFT JOIN jc jcf ON tnr_shortcuts.source = jcf.id
	LEFT JOIN jc jct ON tnr_shortcuts.target = jct.id;

CREATE TABLE tnr_access (
  junction BIGINT,
  access BIGINT[],
  cost DOUBLE PRECISION[]
);

COPY tnr_access(junction, access, cost)
FROM '/docker-entrypoint-initdb.d/preprocessings/tnr_access.csv'
DELIMITER ';'
CSV HEADER;