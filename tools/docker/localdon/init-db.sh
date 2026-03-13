#!/bin/bash
set -e

for db in node1 node2 node3 node4; do
  echo "Creating database: $db"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE $db;
EOSQL
done
