#!/usr/bin/env bash
# dynamic-kubefwd.sh
#
# Automatically manages port-forwarding to Kubernetes services using kubefwd,
# with real-time detection and adaptation to service changes. When services
# are added, removed, or modified in the specified namespace, this script
# automatically restarts the port forwarding to reflect the changes.
#
# Features:
# - Continuous monitoring of Kubernetes services in a namespace
# - Automatic restart of kubefwd when service definitions change
# - Clean process management with proper signal handling
# - Throttled restarts to prevent rapid fluctuations
# - Detailed logging of service changes and forwarding status
#
# Requirements:
# - kubectl configured with access to the target cluster
# - kubefwd installed (https://github.com/txn2/kubefwd)
# - sudo access (required by kubefwd)
#
# Usage: ./dynamic-kubefwd.sh <namespace>
set -eo pipefail

# Usage: ./live-kubefwd.sh <namespace>
if [ -z "$1" ]; then
  echo "Usage: $0 <namespace>"
  exit 1
fi

NAMESPACE="$1"
KUBEFWD_CMD="sudo kubefwd svc -n $NAMESPACE"
LAST_SVCS=""
KPID=""

log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Clean up function
cleanup() {
  log "Cleaning up..."
  [ -n "$KPID" ] && { kill $KPID 2>/dev/null || true; }
  [ -n "$WATCH_PID" ] && { kill $WATCH_PID 2>/dev/null || true; }
  exit 0
}

# Set up signal handlers
trap cleanup INT TERM EXIT

# Function to list services (by name) sorted
list_svcs() {
  kubectl get svc -n "$NAMESPACE" --no-headers 2>/dev/null \
    | awk '{print $1}' \
    | sort \
    | tr '\n' ' ' || echo "Error fetching services"
}

# Function to start kubefwd
start_kubefwd() {
  if [ -n "$KPID" ]; then
    log "Stopping previous kubefwd (PID: $KPID)"
    kill $KPID 2>/dev/null || true
    sleep 1
  fi

  log "Starting kubefwd for namespace '$NAMESPACE'"
  $KUBEFWD_CMD &
  KPID=$!
  sleep 2

  if ! ps -p $KPID > /dev/null; then
    log "Failed to start kubefwd"
    exit 1
  fi

  log "kubefwd started successfully (PID: $KPID)"
}

# Initial launch
LAST_SVCS=$(list_svcs)
start_kubefwd

log "Watching namespace '$NAMESPACE' for services: $LAST_SVCS"
LAST_RESTART=$(date +%s)

# Use polling approach instead of --watch
while true; do
  NEW_SVCS=$(list_svcs)
  CURRENT_TIME=$(date +%s)

  if [[ "$NEW_SVCS" != "$LAST_SVCS" ]]; then
    if (( CURRENT_TIME - LAST_RESTART >= 5 )); then
      log "Service list changed:"
      log "  before: $LAST_SVCS"
      log "   after: $NEW_SVCS"
      start_kubefwd
      LAST_SVCS=$NEW_SVCS
      LAST_RESTART=$CURRENT_TIME
    else
      log "Service change detected, but throttling restarts. Will update shortly."
    fi
  fi

  # Poll every 2 seconds
  sleep 2
done