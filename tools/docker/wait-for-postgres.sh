#!/usr/bin/env bash

RETRIES=5
THRESHOLD=2

until [ $THRESHOLD -eq 0 ] || [ $RETRIES -eq 0 ]; do
  if pg_isready $@; then
    ((THRESHOLD--))
  fi
  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
  sleep 2
done

if [ $THRESHOLD -eq 0 ]; then exit 0; fi
exit 1
