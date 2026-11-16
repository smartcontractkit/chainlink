#!/usr/bin/env bash
#
# capability-mode.sh — choose which cron and consensus each workflow node runs.
#
#   ./capability-mode.sh legacy   every node runs the legacy LOOPs (cron, consensus)
#   ./capability-mode.sh v2       every node runs the v2 binaries (cronv2, consensusv2)
#   ./capability-mode.sh both     the first two nodes run legacy, the other two v2
#
# A node runs v2 when its CRONV2ActivePeriod / ConsensusV2ActivePeriod CRE settings
# are active. This only rewrites those settings, in the CL_CRE_SETTINGS_DEFAULT of
# each node spec marked "# capability-mode" in $TOPO. Takes effect on the next start.
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TOPO="${TOPO:-$SCRIPT_DIR/configs/workflow-gateway-capabilities-don-cron-consensus.toml}"

ACTIVE='[2000-01-01 00:00:00 +0000 UTC,2100-01-01 00:00:00 +0000 UTC]'
INACTIVE='[2100-01-01 00:00:00 +0000 UTC,2101-01-01 00:00:00 +0000 UTC]'

mode="${1:-}"
case "$mode" in
  legacy|v2|both) ;;
  *) echo "usage: $0 legacy|v2|both" >&2; exit 1 ;;
esac

[ -f "$TOPO" ] || { echo "ERROR: topology not found: $TOPO" >&2; exit 1; }

tmp="$TOPO.tmp"
awk -v mode="$mode" -v active="$ACTIVE" -v inactive="$INACTIVE" '
  /# capability-mode$/ {
    node++
    period = inactive
    if (mode == "v2" || (mode == "both" && node > 2)) {
      period = active
    }

    settings = "CL_CRE_SETTINGS_DEFAULT = '\''{\"CRONV2ActivePeriod\":\"" period "\",\"ConsensusV2ActivePeriod\":\"" period "\"}'\''"
    sub(/CL_CRE_SETTINGS_DEFAULT = '\''[^'\'']*'\''/, settings)
    print
    next
  }
  { print }
  END {
    if (node != 4) {
      printf "ERROR: expected 4 node specs marked # capability-mode, found %d\n", node > "/dev/stderr"
      exit 1
    }
  }
' "$TOPO" > "$tmp"
mv "$tmp" "$TOPO"

echo "==> $(basename "$TOPO"): workflow nodes set to $mode"
awk -v active="$ACTIVE" '
  /# capability-mode$/ {
    printf "    workflow-node%d: %s\n", node++, (index($0, active) ? "v2" : "legacy")
  }
' "$TOPO"
