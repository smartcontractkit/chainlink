#!/usr/bin/env bash
#
# Runs the enginefleet container, measures its memory usage over the test, and
# (by default) leaves the container running afterwards so you can inspect memory.
#
# It samples `docker stats` at a fixed interval, writes a CSV timeline, and
# reports the peak/average plus the exact kernel-tracked cgroup peak.
#
# Usage:
#   ./run.sh -w workflow.wasm [-c config.json] [-n instances] [options]
#
# Options:
#   -w PATH   Path to the workflow WASM binary (required, must be UNCOMPRESSED).
#   -c PATH   Path to the workflow config file.
#   -n N      Number of engine instances to start (default 1).
#   -t SECS   End sampling after SECS (default: until the container exits/Ctrl-C).
#   -i SECS   Memory sampling interval in seconds (default 1).
#   -o PATH   CSV output path (default ./mem-<timestamp>.csv).
#   -m LIMIT  Apply a container memory limit (e.g. 512m, 2g).
#   -P PORT   Host port to publish the in-container pprof server on (default 6060).
#   -j        Enable jemalloc native-heap profiling (for CGo/wasmtime memory).
#             Dumps land in ./jeprof-<ts>/ (mounted at /prof); analyse with jeprof.
#   -x        Remove the container when sampling ends (default: keep it running).
#   -s        Enable SuspendOnAwait (executions suspend the WASM instance while
#             awaiting capability responses instead of blocking).
#   -T TAG    Image tag to run (default enginefleet:latest).
#   -d        Enable engine debug logging.
#
set -euo pipefail

IMAGE="enginefleet:latest"
INSTANCES=1
DURATION=""
INTERVAL=1
CSV=""
MEM_LIMIT=""
DEBUG=""
WASM=""
CONFIG=""
PPROF_PORT=6060
JEMALLOC=0
KEEP=1
SUSPEND=""
PROF_DIR=""

while getopts ":w:c:n:t:i:o:m:P:jxsT:dh" opt; do
  case "$opt" in
    w) WASM="$OPTARG" ;;
    c) CONFIG="$OPTARG" ;;
    n) INSTANCES="$OPTARG" ;;
    t) DURATION="$OPTARG" ;;
    i) INTERVAL="$OPTARG" ;;
    o) CSV="$OPTARG" ;;
    m) MEM_LIMIT="$OPTARG" ;;
    P) PPROF_PORT="$OPTARG" ;;
    j) JEMALLOC=1 ;;
    x) KEEP=0 ;;
    s) SUSPEND="-suspend" ;;
    T) IMAGE="$OPTARG" ;;
    d) DEBUG="-debug" ;;
    h) sed -n '2,36p' "$0"; exit 0 ;;
    *) echo "Invalid option. Run with -h for usage." >&2; exit 1 ;;
  esac
done

if [[ -z "$WASM" ]]; then
  echo "-w (path to workflow WASM binary) is required" >&2
  exit 1
fi
if [[ ! -f "$WASM" ]]; then
  echo "workflow binary not found: $WASM" >&2
  exit 1
fi
WASM_ABS="$(cd "$(dirname "$WASM")" && pwd)/$(basename "$WASM")"

NAME="enginefleet-bench-$$"
if [[ -z "$CSV" ]]; then
  CSV="./mem-$(date +%Y%m%d-%H%M%S).csv"
fi

# Assemble docker run arguments.
run_args=(run -d --name "$NAME" -v "$WASM_ABS":/workflow.wasm:ro)
if [[ "$PPROF_PORT" != "0" ]]; then
  run_args+=(-p "${PPROF_PORT}:6060")
fi
app_args=(-n "$INSTANCES" -w /workflow.wasm)

if [[ -n "$CONFIG" ]]; then
  if [[ ! -f "$CONFIG" ]]; then
    echo "config file not found: $CONFIG" >&2
    exit 1
  fi
  CONFIG_ABS="$(cd "$(dirname "$CONFIG")" && pwd)/$(basename "$CONFIG")"
  run_args+=(-v "$CONFIG_ABS":/config:ro)
  app_args+=(-c /config)
fi
if [[ -n "$MEM_LIMIT" ]]; then
  run_args+=(--memory "$MEM_LIMIT")
fi
if [[ "$JEMALLOC" == "1" ]]; then
  PROF_DIR="./jeprof-$(date +%Y%m%d-%H%M%S)"
  mkdir -p "$PROF_DIR"
  PROF_ABS="$(cd "$PROF_DIR" && pwd)"
  run_args+=(-v "$PROF_ABS":/prof)
  run_args+=(-e "LD_PRELOAD=libjemalloc.so.2")
  # Sample every ~512KiB (lg_prof_sample:19), auto-dump each ~1GiB allocated
  # (lg_prof_interval:30), and dump once more on clean exit (prof_final).
  run_args+=(-e "MALLOC_CONF=prof:true,prof_active:true,prof_prefix:/prof/jeprof,lg_prof_sample:19,lg_prof_interval:30,prof_final:true")
fi
if [[ -n "$SUSPEND" ]]; then
  app_args+=("$SUSPEND")
fi
if [[ -n "$DEBUG" ]]; then
  app_args+=("$DEBUG")
fi

# Convert a docker stats memory string (e.g. "45.63MiB") to bytes.
to_bytes() {
  local v="$1" num unit factor
  num="$(sed -E 's/([0-9.]+).*/\1/' <<<"$v")"
  unit="$(sed -E 's/[0-9.]+//' <<<"$v")"
  case "$unit" in
    B)   factor=1 ;;
    KiB) factor=1024 ;;
    MiB) factor=1048576 ;;
    GiB) factor=1073741824 ;;
    TiB) factor=1099511627776 ;;
    kB)  factor=1000 ;;
    MB)  factor=1000000 ;;
    GB)  factor=1000000000 ;;
    *)   factor=1 ;;
  esac
  awk -v n="$num" -v f="$factor" 'BEGIN { printf "%.0f", n * f }'
}

human() { awk -v b="$1" 'BEGIN { printf "%.2f MiB", b / 1048576 }'; }

container_running() {
  [[ "$(docker inspect -f '{{.State.Running}}' "$NAME" 2>/dev/null)" == "true" ]]
}

STOP=0
trap 'echo; echo ">> Interrupted"; STOP=1' INT TERM

echo ">> Starting container ${NAME} (image ${IMAGE}, instances ${INSTANCES})"
docker "${run_args[@]}" "$IMAGE" "${app_args[@]}" >/dev/null

echo "timestamp_unix,mem_bytes,mem_human" > "$CSV"

peak=0
sum=0
count=0
cgroup_peak=0
start_ts="$(date +%s)"

echo ">> Sampling memory every ${INTERVAL}s (CSV: ${CSV})"
while true; do
  container_running || { echo ">> Container exited"; break; }

  raw="$(docker stats --no-stream --format '{{.MemUsage}}' "$NAME" 2>/dev/null || true)"
  used="${raw%% /*}"
  used="${used// /}"
  if [[ -n "$used" ]]; then
    bytes="$(to_bytes "$used")"
    now="$(date +%s)"
    echo "${now},${bytes},${used}" >> "$CSV"
    (( bytes > peak )) && peak="$bytes"
    sum=$(( sum + bytes ))
    count=$(( count + 1 ))
  fi

  # Best-effort exact peak from the kernel (cgroup v2).
  cg="$(docker exec "$NAME" cat /sys/fs/cgroup/memory.peak 2>/dev/null || true)"
  if [[ "$cg" =~ ^[0-9]+$ ]] && (( cg > cgroup_peak )); then
    cgroup_peak="$cg"
  fi

  if [[ "$STOP" == "1" ]]; then break; fi
  if [[ -n "$DURATION" ]]; then
    now="$(date +%s)"
    if (( now - start_ts >= DURATION )); then
      echo ">> Duration ${DURATION}s reached, ending sampling"
      break
    fi
  fi
  sleep "$INTERVAL"
done

# Final cgroup peak before we report.
cg="$(docker exec "$NAME" cat /sys/fs/cgroup/memory.peak 2>/dev/null || true)"
if [[ "$cg" =~ ^[0-9]+$ ]] && (( cg > cgroup_peak )); then
  cgroup_peak="$cg"
fi

elapsed=$(( $(date +%s) - start_ts ))

echo
echo "================ memory report ================"
echo "container:        ${NAME}"
echo "image:            ${IMAGE}"
echo "instances:        ${INSTANCES}"
echo "duration:         ${elapsed}s"
echo "samples:          ${count}"
echo "sampled peak:     $(human "$peak")"
if (( count > 0 )); then
  echo "sampled average:  $(human $(( sum / count )))"
fi
if (( cgroup_peak > 0 )); then
  echo "cgroup peak:      $(human "$cgroup_peak")   (exact, kernel-tracked)"
fi
echo "csv:              ${CSV}"
echo "==============================================="

if [[ "$KEEP" == "0" ]]; then
  echo ">> Removing container"
  docker rm -f "$NAME" >/dev/null 2>&1 || true
  exit 0
fi

cat <<EOF

>> Container '${NAME}' is still running for inspection. Memory inspection:

  # live totals
  docker stats ${NAME}
  docker exec ${NAME} cat /sys/fs/cgroup/memory.current   # current bytes
  docker exec ${NAME} cat /sys/fs/cgroup/memory.peak      # high-water bytes

  # Go heap (managed memory)
  go tool pprof http://localhost:${PPROF_PORT}/debug/pprof/heap

  # native/mmap regions (wasmtime linear memory, compiled code) — the C side
  docker exec ${NAME} cat /proc/1/smaps_rollup
  docker exec ${NAME} pmap -x 1 | sort -k3 -n | tail -20
EOF

if [[ "$JEMALLOC" == "1" ]]; then
  cat <<EOF

  # jemalloc native heap profile (CGo/wasmtime malloc); dumps in ${PROF_DIR}/
  ls -t ${PROF_DIR}/*.heap
  docker exec ${NAME} sh -c 'jeprof --show_bytes --text /usr/local/bin/enginefleet \$(ls -t /prof/*.heap | head -1)'
EOF
fi

echo
echo ">> Stop it when done:  docker rm -f ${NAME}"
