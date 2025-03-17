#!/bin/bash

set -xe

export PGCLIENTENCODING=latin1

# Import teleatlas dataset
ogr2ogr -nln nw -nlt PROMOTE_TO_MULTI -lco GEOMETRY_NAME=geom -lco FID=gid -lco PRECISION=NO PG:"dbname=$POSTGRES_DB user=$POSTGRES_USER port=5432" /tmp/teleatlas/nw.shp
ogr2ogr -nln gc -nlt PROMOTE_TO_MULTI -lco GEOMETRY_NAME=geom -lco FID=gid -lco PRECISION=NO PG:"dbname=$POSTGRES_DB user=$POSTGRES_USER port=5432" /tmp/teleatlas/gc.shp
ogr2ogr -nln jc -nlt PROMOTE_TO_MULTI -lco GEOMETRY_NAME=geom -lco FID=gid -lco PRECISION=NO PG:"dbname=$POSTGRES_DB user=$POSTGRES_USER port=5432" /tmp/teleatlas/jc.shp
ogr2ogr -nln nl -nlt PROMOTE_TO_MULTI -lco GEOMETRY_NAME=geom -lco FID=gid -lco PRECISION=NO PG:"dbname=$POSTGRES_DB user=$POSTGRES_USER port=5432" /tmp/teleatlas/nl.dbf
