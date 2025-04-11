#!/bin/bash

set -xe

pg_restore -U "$POSTGRES_USER" -d "$POSTGRES_DB" /docker-entrypoint-initdb.d/postgres_localhost-2025_04_11_10_19_36-dump.sql
