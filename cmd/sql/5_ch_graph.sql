CREATE TABLE ch_junctions (
  junction BIGINT,
  orderpos BIGINT,
  importance BIGINT
);

COPY ch_junctions(junction, orderpos, importance)
FROM '/docker-entrypoint-initdb.d/preprocessings/ch_junctions.csv'
DELIMITER ';'
CSV HEADER;

CREATE TABLE ch_shortcuts (
  source BIGINT,
  target BIGINT,
  cost DOUBLE PRECISION,
  via BIGINT
);

COPY ch_shortcuts(source, target, cost, via)
FROM '/docker-entrypoint-initdb.d/preprocessings/ch_shortcuts.csv'
DELIMITER ';'
CSV HEADER;

CREATE TABLE ch_graph (
  source BIGINT,
  target BIGINT,
  cost DOUBLE PRECISION
);

COPY tnr_access(source, target, cost)
FROM '/docker-entrypoint-initdb.d/preprocessings/ch_graph.csv'
DELIMITER ';'
CSV HEADER;
