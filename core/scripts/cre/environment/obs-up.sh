#!/usr/bin/env bash
# obs-up.sh — wrapper around `./bin/ctf obs up` that applies local overrides.
#
# CTF regenerates the compose/ directory on every `obs up`, overwriting any
# manual edits.  This script re-applies our customisations afterwards and
# restarts only the affected containers so Grafana and the OTel collector pick
# them up without a full stack restart.
#
# Usage:
#   ./obs-up.sh             # bring the stack up (or recreate it) with patches
#   ./obs-up.sh --down      # tear down first, then bring up with patches
#
# Tracked overrides live in observability-overrides/ and are applied here.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OVERRIDES="$SCRIPT_DIR/observability-overrides"
COMPOSE_DIR="$SCRIPT_DIR/compose"

# ── Optionally tear down first ────────────────────────────────────────────────
if [[ "${1:-}" == "--down" ]]; then
  echo "▼  Tearing down obs stack..."
  "$SCRIPT_DIR/bin/ctf" obs down || true
fi

# ── Bring up the stack (CTF regenerates compose/ here) ───────────────────────
echo "▲  Starting obs stack..."
"$SCRIPT_DIR/bin/ctf" obs up

# ── Apply otel-collector override ─────────────────────────────────────────────
echo "⚙  Applying otel.yaml override (resource_to_telemetry_conversion)..."
cp "$OVERRIDES/otel.yaml" "$COMPOSE_DIR/otel.yaml"

# ── Apply dashboard ───────────────────────────────────────────────────────────
echo "⚙  Copying Durable Emitter Load Test dashboard..."
mkdir -p "$COMPOSE_DIR/conf/provisioning/dashboards/beholder"
cp "$OVERRIDES/dashboards/beholder/load_dashboard.json" \
   "$COMPOSE_DIR/conf/provisioning/dashboards/beholder/load_dashboard.json"

# ── Patch docker-compose.yaml to add beholder dashboard volume mount ──────────
echo "⚙  Patching docker-compose.yaml to add dashboard volume mount..."
python3 - "$COMPOSE_DIR/docker-compose.yaml" <<'PYEOF'
import sys, re

path = sys.argv[1]
with open(path) as f:
    content = f.read()

marker = "./conf/provisioning/dashboards/beholder/load_dashboard.json:/var/lib/grafana/dashboards/beholder/load_dashboard.json"
if marker in content:
    print("   dashboard volume mount already present, skipping.")
    sys.exit(0)

# Insert our mount after the last existing dashboard volume line.
content = re.sub(
    r"([ \t]+- \./conf/provisioning/dashboards/workflow-engine/engine\.json:/var/lib/grafana/dashboards/workflow-engine/engine\.json)",
    r"\1\n      - ./conf/provisioning/dashboards/beholder/load_dashboard.json:/var/lib/grafana/dashboards/beholder/load_dashboard.json",
    content,
)

with open(path, "w") as f:
    f.write(content)
print("   done.")
PYEOF

# ── Recreate affected containers so new volume mounts and config take effect ──
# `restart` reuses the existing container spec (no new mounts); `up --force-recreate`
# rebuilds the container from the patched docker-compose.yaml.
echo "↺  Recreating otel-collector and grafana with updated config..."
docker compose -f "$COMPOSE_DIR/docker-compose.yaml" up -d --force-recreate otel-collector grafana

echo ""
echo "✓  Obs stack is up with all overrides applied."
echo "   Grafana:        http://localhost:3000"
echo "   Dashboard:      http://localhost:3000/d/durable-emitter-load-test"
