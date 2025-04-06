#!/bin/sh

set -xe

docker-compose -f docker/compose.yaml --env-file .env up -d --build
