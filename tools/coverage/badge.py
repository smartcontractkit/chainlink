#!/usr/bin/env python3
"""Generate a small coverage badge and latest-report JSON from tools/coverage output.

Usage:
  badge.py any.json badge.svg report.json SHA artifact_name repo run_url

The badge uses the source-only global percentage, which excludes tests,
generated code, and mocks. The report JSON is the only state kept in the
orphan `coverage` branch.
"""

from __future__ import annotations

import datetime as dt
import json
import sys
from pathlib import Path


def color_for(percent: float) -> str:
    if percent >= 70:
        return "#2e7d32"
    if percent >= 55:
        return "#f9a825"
    if percent >= 45:
        return "#ef6c00"
    return "#c62828"


def main() -> int:
    if len(sys.argv) != 8:
        print(__doc__, file=sys.stderr)
        return 2

    report_path, svg_path, state_path, sha, artifact, repo, run_url = sys.argv[1:]
    report = json.loads(Path(report_path).read_text(encoding="utf-8"))
    source = report.get("GlobalSourceOnly") or report.get("Global") or {}
    percent = float(source.get("Percent", 0.0))
    raw = float((report.get("Global") or {}).get("Percent", 0.0))
    patch = report.get("Patch") or {}
    patch_percent = patch.get("Percent")

    value = f"{percent:.1f}%"
    color = color_for(percent)
    svg = (
        '<svg xmlns="http://www.w3.org/2000/svg" width="160" height="20" '
        f'role="img" aria-label="coverage: {value}">'
        '<rect width="160" height="20" rx="3" fill="#555"/>'
        f'<rect x="80" width="80" height="20" rx="3" fill="{color}"/>'
        f'<rect x="80" width="4" height="20" fill="{color}"/>'
        '<g fill="#fff" text-anchor="middle" '
        'font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11">'
        '<text x="40" y="14">coverage</text>'
        f'<text x="120" y="14">{value}</text>'
        "</g></svg>\n"
    )
    Path(svg_path).write_text(svg, encoding="utf-8")

    state = {
        "sha": sha,
        "repo": repo,
        "global_raw": raw,
        "global_source_only": percent,
        "patch": patch_percent,
        "report_artifact": artifact,
        "run_url": run_url,
        "updated_at": dt.datetime.now(dt.timezone.utc)
        .replace(microsecond=0)
        .isoformat()
        .replace("+00:00", "Z"),
    }
    Path(state_path).write_text(json.dumps(state, indent=2) + "\n", encoding="utf-8")
    print(f"coverage badge: {value} ({color})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
