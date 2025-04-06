#!/bin/bash
set -e
pg_restore -U "$POSTGRES_USER" -d "$POSTGRES_DB" /docker-entrypoint-initdb.d/postgres_localhost-2025_04_06_15_26_57-dump.sql
