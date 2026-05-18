#!/usr/bin/env bash

DB_FILE="$1"
DB_SUPER_USER="postgres"
DB_USER="postgres"
DB_NAME="chainlink_fallback_db"
DB_HOST_PORT="localhost:5432"

psql "postgresql://$DB_SUPER_USER@$DB_HOST_PORT/postgres" -c "CREATE DATABASE $DB_NAME"
psql "postgresql://$DB_SUPER_USER@$DB_HOST_PORT/postgres" -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

pg_restore -d "postgresql://$DB_SUPER_USER@$DB_HOST_PORT/$DB_NAME" "$DB_FILE"