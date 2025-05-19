#!/usr/bin/env bash

while [ "$#" -gt 1 ] && [ "$1" != "--" ]; do
    /opt/docker-wait-for-it.sh $1
    shift
done

# Hand off to the CMD
exec "$@"
