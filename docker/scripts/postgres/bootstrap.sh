#!/bin/bash

set -xe

export PGCLIENTENCODING=latin1

# Import teleatlas dataset
for shapefile in /tmp/teleatlas/*.shp; do
    tablename=$(basename $shapefile .shp)
    ogr2ogr -nln $tablename -nlt PROMOTE_TO_MULTI -lco GEOMETRY_NAME=geom -lco FID=gid -lco PRECISION=NO PG:"dbname=$POSTGRES_DB user=$POSTGRES_USER port=5432" $shapefile
done
