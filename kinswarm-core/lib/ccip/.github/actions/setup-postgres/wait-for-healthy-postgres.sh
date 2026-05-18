#!/bin/bash
RETRIES=10

until [ $RETRIES -eq 0 ]; do
  DOCKER_OUTPUT=$(docker compose ps postgres --status running --format json)
  JSON_TYPE=$(echo "$DOCKER_OUTPUT" | jq -r 'type')

  if [ "$JSON_TYPE" == "array" ]; then
    HEALTH_STATUS=$(echo "$DOCKER_OUTPUT" | jq -r '.[0].Health')
  elif [ "$JSON_TYPE" == "object" ]; then
    HEALTH_STATUS=$(echo "$DOCKER_OUTPUT" | jq -r '.Health')
  else
    HEALTH_STATUS="Unknown JSON type: $JSON_TYPE"
  fi

  echo "postgres health status: $HEALTH_STATUS"
  if [ "$HEALTH_STATUS" == "healthy" ]; then
    exit 0
  fi

  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
  sleep 2
done

exit 1
