#!/bin/bash
set -e
set -x

all_files="-f docker-compose.yaml -f docker-compose.gethnet.yaml -f docker-compose.postgres.yaml -f docker-compose.integration.yaml"
docker-compose $all_files build --parallel
docker-compose $all_files up --exit-code-from integration