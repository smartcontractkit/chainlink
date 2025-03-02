#!/usr/bin/env bash
# Description: Waits until the number of running Docker containers equals the target count,
# or until a specified timeout is reached.
#
# Usage:
#   ./wait-for-docker-containers.sh [timeout_seconds] [target_container_count]
#
#   timeout_seconds         - Optional: Maximum seconds to wait (default: 30)
#   target_container_count  - Optional: Container count to wait for (default: 0)

# Read parameters or use default values.
TIMEOUT=${1:-30}
TARGET_COUNT=${2:-0}

# Calculate the end time using the built-in SECONDS variable.
end=$((SECONDS + TIMEOUT))

# Loop until the current docker container count equals the target count.
while [ "$(docker ps -q | wc -l)" -ne "$TARGET_COUNT" ]; do
  # If the timeout has been reached, exit with an error.
  if [ $SECONDS -ge $end ]; then
    exit 1
  fi
  sleep 1
done