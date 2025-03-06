#!/bin/sh

set -xe

docker-compose -f docker/compose.yaml --env-file .env up -d --build

{
    echo "Neo4j running on: http://localhost:7474"
    echo "PostgreSQL running on: http://localhost:5432"
    echo "Server running on: http://localhost:8080"
} 2> /dev/null

