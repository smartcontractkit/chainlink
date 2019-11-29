#!/bin/bash
set -e

all_files="-f docker-compose.yaml -f docker-compose.paritynet.yaml -f docker-compose.postgres.yaml -f docker-compose.integration.yaml"
base="docker-compose $all_files"
dev="$base -f docker-compose.dev.yaml"
usage="parity_postgres_test -- run the integration test suite against parity as the ethereum node, and postgres as the backing database for chainlink\

Commands:
    base        Echo the docker-compose command with all of its composed files,
                useful for running docker-compose commands that aren't shortcutted here
    down        Brings down all services.
    down:clean  Brings down all services and removes any volumes, good for clean slate testing.
    up          Brings up all services.
    up:dev      Brings up all services, and bind-mounts source files for quick changes without rebuilding a container.
                See docker-compose.dev.yaml for the list of bind-mounted folders per service.
    build       Builds all images in parallel.
    test        Runs the test suite, exiting on any failures."

case "$1" in
  help)
    echo "$usage"
    ;;
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
