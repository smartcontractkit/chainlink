#!/usr/bin/env bash
# obs-up.sh — wrapper around `./bin/ctf obs up` that applies local overrides for durable-emitter load tests.
# Usage:
#   ./obs-up.sh             # bring the stack up (or recreate it) with patches
#   ./obs-up.sh --down      # tear down first, then bring up with patches
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

# ── Patch docker-compose.yaml to add beholder dashboard volume mount,
#    grafana-image-renderer service, and Grafana rendering env vars ─────────────
echo "⚙  Patching docker-compose.yaml (dashboard mount + image renderer)..."
python3 - "$COMPOSE_DIR/docker-compose.yaml" <<'PYEOF'
import sys

path = sys.argv[1]
with open(path) as f:
    lines = f.readlines()

out = []
i = 0
added_dashboard = False
added_env = False
added_renderer_dep = False
added_renderer_svc = False

while i < len(lines):
    line = lines[i]

    # 1. Dashboard volume mount: insert after the workflow-engine volume line
    if (not added_dashboard
        and "workflow-engine/engine.json:/var/lib/grafana/dashboards/workflow-engine/engine.json" in line
        and "beholder/load_dashboard.json" not in "".join(lines)):
        out.append(line)
        out.append("      - ./conf/provisioning/dashboards/beholder/load_dashboard.json:/var/lib/grafana/dashboards/beholder/load_dashboard.json\n")
        added_dashboard = True
        print("   added dashboard volume mount.")
        i += 1
        continue

    # 2. Grafana env vars: insert after the ports block (line with '3000:3000')
    if (not added_env
        and "'3000:3000'" in line
        and "GF_RENDERING_SERVER_URL" not in "".join(lines)):
        out.append(line)
        out.append("    environment:\n")
        out.append("      GF_RENDERING_SERVER_URL: http://grafana-image-renderer:8081/render\n")
        out.append("      GF_RENDERING_CALLBACK_URL: http://grafana:3000/\n")
        out.append("      GF_LOG_FILTERS: rendering:debug\n")
        added_env = True
        print("   added Grafana rendering env vars.")
        i += 1
        continue

    # 3. Add grafana-image-renderer to depends_on (after '- tempo' in grafana block)
    if (not added_renderer_dep
        and line.strip() == "- tempo"
        and i > 0 and "depends_on" in "".join(lines[max(0,i-3):i])
        and "grafana-image-renderer" not in "".join(lines)):
        out.append(line)
        out.append("      - grafana-image-renderer\n")
        added_renderer_dep = True
        print("   added grafana-image-renderer to depends_on.")
        i += 1
        continue

    # 4. Add grafana-image-renderer service: insert before pyroscope
    if (not added_renderer_svc
        and line.strip() == "pyroscope:"
        and "grafana-image-renderer:" not in "".join(lines)):
        out.append("  grafana-image-renderer:\n")
        out.append("    image: grafana/grafana-image-renderer:latest\n")
        out.append("    ports:\n")
        out.append('      - "8081:8081"\n')
        out.append("    environment:\n")
        out.append('      ENABLE_METRICS: "true"\n')
        out.append("      RENDERING_ARGS: --no-sandbox,--disable-gpu\n")
        out.append('      HTTP_PORT: "8081"\n')
        out.append("\n")
        added_renderer_svc = True
        print("   added grafana-image-renderer service.")

    out.append(line)
    i += 1

with open(path, "w") as f:
    f.writelines(out)

if not added_dashboard:
    print("   dashboard volume mount already present.")
if not added_env:
    print("   Grafana rendering env vars already present.")
if not added_renderer_dep:
    print("   grafana-image-renderer dep already present.")
if not added_renderer_svc:
    print("   grafana-image-renderer service already present.")
print("   done.")
PYEOF

# ── Recreate affected containers so new volume mounts and config take effect ──
# `restart` reuses the existing container spec (no new mounts); `up --force-recreate`
# rebuilds the container from the patched docker-compose.yaml.
echo "↺  Recreating otel-collector, grafana, and image-renderer with updated config..."
docker compose -f "$COMPOSE_DIR/docker-compose.yaml" up -d --force-recreate otel-collector grafana grafana-image-renderer

echo ""
echo "✓  Obs stack is up with all overrides applied."
echo "   Grafana:        http://localhost:3000"
echo "   Dashboard:      http://localhost:3000/d/durable-emitter-load-test"
