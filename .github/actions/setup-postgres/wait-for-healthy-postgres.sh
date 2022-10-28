#!/bin/bash
RETRIES=10

until [ $RETRIES -eq 0 ]; do
  if docker compose ps postgres --status running --format json | jq >/dev/null -e 'if (.[0].Health == "healthy") then true else false end'; then
    exit 0
  fi
  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
  sleep 2
done
exit 1
