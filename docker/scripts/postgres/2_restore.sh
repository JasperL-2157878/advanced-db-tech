#!/bin/bash

set -xe

pg_restore -U "$POSTGRES_USER" -d "$POSTGRES_DB" /docker-entrypoint-initdb.d/postgres_localhost-2025_04_06_16_10_06-dump.sql
