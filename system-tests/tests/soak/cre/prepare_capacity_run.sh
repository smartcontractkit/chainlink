#!/usr/bin/env bash
# Prepares a clean run of Test_V2_CRE_WorkflowCapacity (workflow_capacity_test.go). Run it before every run:
#   1. stops local CRE (`env stop`), so the test starts a fresh environment
#   2. deletes workflow artifacts left by earlier capacity runs (GBs at thousands of workflows)
#   3. recreates the observability stack, wiping Prometheus data and the previous environment's series
#      (RESET_OBSERVABILITY=0 keeps the stack and its history, e.g. to look at the last run in Grafana)
#   4. adjusts the OTel collector config, which `obs up` rewrites every time:
#      - raises the otlp grpc receive limit: with the default 4 MiB, node metric exports are rejected
#        from ~600 workflows on
#      - drops per-workflow limiter metrics (*PerWorkflow*, one series set per workflow): at 3000
#        workflows they were 540k of 580k series, Prometheus scrapes of the collector timed out and
#        node metrics read as missing. Nodes still produce and export them; only the collector drops them.
#
# Usage: system-tests/tests/soak/cre/prepare_capacity_run.sh
# Env:   OTEL_MAX_RECV_MSG_SIZE_MIB (default 128), RESET_OBSERVABILITY (default 1)
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/../../../.." && pwd)
ENV_DIR="$REPO_ROOT/core/scripts/cre/environment"
ARTIFACTS_DIR="${TMPDIR:-/tmp}/cre-workflow-artifacts"
OTEL_CONFIG="$HOME/.local/share/ctf/compose/otel.yaml"
OTEL_CONTAINER=compose-otel-collector-1
PROMETHEUS_CONTAINER=compose-prometheus-1
export MAX_RECV_MIB=${OTEL_MAX_RECV_MSG_SIZE_MIB:-128}
RESET_OBSERVABILITY=${RESET_OBSERVABILITY:-1}

log() { printf '\n==> %s\n' "$*"; }
container_running() { [ "$(docker inspect -f '{{.State.Running}}' "$1" 2>/dev/null)" = "true" ]; }

log "Stopping local CRE"
(cd "$ENV_DIR" && go run . env stop)
leftover=$(docker ps -aq --filter "label=framework=ctf" | wc -l | tr -d ' ')
if [ "$leftover" != "0" ]; then
	echo "error: $leftover local CRE containers still exist; remove them with: docker ps -aq --filter label=framework=ctf | xargs docker rm -f" >&2
	exit 1
fi

log "Removing artifacts of earlier capacity runs from $ARTIFACTS_DIR"
rm -rf "$ARTIFACTS_DIR"/test_v2_cre_workflowcapacity-*

if [ "$RESET_OBSERVABILITY" = "1" ]; then
	log "Recreating the observability stack (wipes Prometheus data)"
	(cd "$ENV_DIR" && go run . obs down && go run . obs up -f)
elif ! container_running "$PROMETHEUS_CONTAINER" || ! container_running "$OTEL_CONTAINER"; then
	log "Starting the observability stack"
	(cd "$ENV_DIR" && go run . obs up -f)
fi

log "Configuring the OTel collector (receive limit ${MAX_RECV_MIB} MiB, drop *PerWorkflow* metrics)"
if [ ! -f "$OTEL_CONFIG" ]; then
	echo "error: $OTEL_CONFIG not found; is the observability stack from \`go run . obs up -f\`?" >&2
	exit 1
fi
perl -0pi -e '
	# receive limit, right below the grpc endpoint with the same indentation
	s/^[ \t]*max_recv_msg_size_mib:.*\n//mg;
	s/^([ \t]*)(endpoint: 0\.0\.0\.0:4317\n)/$1$2$1max_recv_msg_size_mib: $ENV{MAX_RECV_MIB}\n/m;
	# filter processor, used by the metrics pipeline only
	if (!/filter\/perworkflow:/) {
		s/^processors:\n/processors:\n  filter\/perworkflow:\n    error_mode: ignore\n    metrics:\n      metric:\n        - '"'"'IsMatch(name, ".*PerWorkflow.*")'"'"'\n/m;
		s/(\n[ \t]*metrics:\n[ \t]*receivers: \[otlp\]\n[ \t]*processors: \[)batch\]/$1filter\/perworkflow, batch]/;
	}
' "$OTEL_CONFIG"
if ! grep -q "max_recv_msg_size_mib: $MAX_RECV_MIB" "$OTEL_CONFIG" || ! grep -q "processors: \[filter/perworkflow, batch\]" "$OTEL_CONFIG"; then
	echo "error: could not apply the collector changes to $OTEL_CONFIG; its layout may have changed" >&2
	exit 1
fi

docker restart "$OTEL_CONTAINER" >/dev/null
sleep 5
if ! container_running "$OTEL_CONTAINER"; then
	echo "error: $OTEL_CONTAINER did not come back after the config change:" >&2
	docker logs --tail 20 "$OTEL_CONTAINER" >&2
	exit 1
fi

log "Ready"
docker info --format 'Docker memory: {{.MemTotal}} bytes, CPUs: {{.NCPU}}'
cat <<EOF

Run the test from system-tests/tests, e.g. 5000 workflows on the 2-node workflow DON:

  CTF_CONFIGS=configs/workflow-gateway-don-capacity-test-2-nodes.toml \\
  CRE_CAPACITY_TARGET=5000 CRE_CAPACITY_STEP=500 \\
    go test -v -count=1 -timeout 6h -run '^Test_V2_CRE_WorkflowCapacity\$' ./soak/cre 2>&1 | tee run-5000.txt
EOF
