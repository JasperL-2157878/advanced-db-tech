#!/bin/sh

set -xe

if [ "$#" -ne 2 ]; then
    echo "Usage: import_shapefile.sh <shapefile> <table_name>"
    exit 1
fi

export PGCLIENTENCODING=latin1
ogr2ogr -nln $2 -nlt PROMOTE_TO_MULTI -lco GEOMETRY_NAME=geom -lco FID=gid -lco PRECISION=NO Pg:"dbname=postgres host=localhost user=postgres password=postgres port=5432" $1