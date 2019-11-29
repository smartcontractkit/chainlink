#!/bin/bash
set -e

all_files="-f docker-compose.yaml -f docker-compose.paritynet.yaml -f docker-compose.postgres.yaml -f docker-compose.integration.yaml"
base="docker-compose $all_files"
dev="$base -f docker-compose.dev.yaml"

case "$1" in
  base)
    echo $base
    ;;
  down)
    $base down
    ;;
  down:clean)
    $base down -v
    ;;
  up)
    $base up
    ;;
  up:dev)
    $dev up
    ;;
  build)
    $base build --parallel
    ;;
  test)
    $base up --exit-code-from integration
    ;;
  *)
    ./parity_postgres_test.sh build
    ./parity_postgres_test.sh test
    ;;
esac
