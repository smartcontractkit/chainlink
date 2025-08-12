#!/usr/bin/env bash

echo "Cleaning up java-tron container.."

echo "Checking for existing 'chainlink-tron.java-tron' docker containers..."

dpid=`docker ps -a | grep chainlink-tron.java-tron | awk '{print $1}'`

if [ -z "$dpid" ]; then
  echo "No docker java-tron containers running."
else
  for id in $dpid; do
    echo "Killing docker container: ${id}"
    docker kill $id || true
    # Try to kill the docker container normally, if it doesn't work, use the --force flag to send a SIGKILL.
    # docker rm --force always exits with code 0, even if the container no longer exists.
    docker rm $id || docker rm --force $id
  done
fi

echo "Removing network.."
docker network rm --force chainlink-tron.network

echo "Cleanup finished."
