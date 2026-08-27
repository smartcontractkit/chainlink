#!/usr/bin/env python3
import argparse
import collections
import json
import os
import sys
from typing import Any, Dict, List, Optional, Sequence, Tuple


def get_artifact_suffix(module_name: str) -> str:
    clean = module_name.strip().rstrip("/")
    if clean.startswith("./"):
        clean = clean[2:]
    return clean.replace("/", "-") if clean and clean != "." else "root"


def process(
    report_path: str,
    module_name: str = "",
    summary_path: Optional[str] = None,
    output_path: Optional[str] = None,
    critical_severities: Sequence[str] = ("high", "medium"),
) -> Tuple[str, Dict[str, Any]]:
    suffix = get_artifact_suffix(module_name)

    if report_path and os.path.isfile(report_path):
        with open(report_path, "r", encoding="utf-8") as f:
            data = json.load(f)
    else:
        print(
            f"::warning::Lint report not found at {report_path!r}; treating as empty report",
            file=sys.stderr,
        )
        data = {"Issues": []}

    issues: List[Dict[str, Any]] = data.get("Issues") or []

    # 1. Metrics & Counts
    filenames = {
        issue.get("Pos", {}).get("Filename")
        for issue in issues
        if issue.get("Pos", {}).get("Filename")
    }
    file_count = len(filenames)
    issue_count = len(issues)

    critical_set = {s.lower() for s in critical_severities}
    critical_issue_count = sum(
        1
        for issue in issues
        if (issue.get("Severity") or "").lower() in critical_set
    )

    linter_counts = collections.Counter(
        issue.get("FromLinter") or "unknown" for issue in issues
    )
    # Sort by count desc, then linter name asc
    sorted_linters = sorted(linter_counts.items(), key=lambda x: (-x[1], x[0]))
    per_source_lines = [f"{count:>4} {linter}" for linter, count in sorted_linters]
    per_source_str = "\n".join(per_source_lines)

    severity_order = {"high": 1, "medium": 2, "low": 3, "warn": 4, "none": 5}
    severity_counts = collections.Counter(
        (issue.get("Severity") or "none").lower() for issue in issues
    )
    sorted_severities = sorted(
        severity_counts.items(),
        key=lambda x: (severity_order.get(x[0], 99), -x[1]),
    )
    per_severity_lines = [f"{count:>4} {sev}" for sev, count in sorted_severities]
    per_severity_str = "\n".join(per_severity_lines)

    outputs = {
        "file-count": file_count,
        "issue-count": issue_count,
        "critical-issue-count": critical_issue_count,
        "per-source-count": per_source_str,
        "per-severity-count": per_severity_str,
        "suffix": suffix,
    }

    # 2. Markdown Step Summary
    display_module_name = module_name.strip()
    if display_module_name in ("", ".", "./"):
        display_module_name = "root"
    summary_lines = [
        f"### 🔍 GolangCI Lint Results: `{display_module_name}`",
        "",
        "| Metric | Count |",
        "| :--- | :--- |",
        f"| **Total Issues** | {issue_count} |",
        f"| **Files Affected** | {file_count} |",
        f"| **Critical Issues** | {critical_issue_count} |",
        "",
    ]

    if issues:
        summary_lines.extend(
            [
                "<details><summary><b>Detailed Issues List</b></summary>",
                "",
                "| File | Line | Linter | Message |",
                "| :--- | :--- | :--- | :--- |",
            ]
        )
        for issue in issues[:100]:
            pos = issue.get("Pos") or {}
            filename = pos.get("Filename", "")
            line = pos.get("Line", "")
            linter = issue.get("FromLinter", "")
            msg = (
                issue.get("Text", "")
                .strip()
                .replace("\r", "")
                .replace("\n", " ")
                .replace("|", "&#124;")
            )
            summary_lines.append(f"| `{filename}` | {line} | `{linter}` | {msg} |")

        summary_lines.append("")
        if len(issues) > 100:
            summary_lines.append(
                "*...showing first 100 issues only (see full report artifact for complete list)*\n"
            )
        summary_lines.append("</details>\n")

    summary_content = "\n".join(summary_lines)

    # Write summary if path provided
    if summary_path:
        with open(summary_path, "a", encoding="utf-8") as f:
            f.write(summary_content + "\n")

    # Write GITHUB_OUTPUT if path provided
    if output_path:
        with open(output_path, "a", encoding="utf-8") as f:
            f.write(f"file-count={file_count}\n")
            f.write(f"issue-count={issue_count}\n")
            f.write(f"critical-issue-count={critical_issue_count}\n")
            f.write(f"suffix={suffix}\n")
            f.write("per-source-count<<EOF\n")
            f.write(per_source_str + "\n")
            f.write("EOF\n")
            f.write("per-severity-count<<EOF\n")
            f.write(per_severity_str + "\n")
            f.write("EOF\n")

    return summary_content, outputs


def main() -> None:
    parser = argparse.ArgumentParser(description="Process GolangCI-Lint JSON report.")
    parser.add_argument(
        "--report-path",
        default=os.environ.get("REPORT_PATH_JSON"),
        help="Path to JSON report file",
    )
    parser.add_argument(
        "--module-name",
        default=os.environ.get("MODULE_NAME", ""),
        help="Module name",
    )
    parser.add_argument(
        "--summary-path",
        default=os.environ.get("GITHUB_STEP_SUMMARY"),
        help="Path to GitHub Step Summary file",
    )
    parser.add_argument(
        "--output-path",
        default=os.environ.get("GITHUB_OUTPUT"),
        help="Path to GitHub Output file",
    )
    critical_default = (
        os.environ["CRITICAL_SEVERITIES"].split(",")
        if "CRITICAL_SEVERITIES" in os.environ
        else ["high", "medium"]
    )
    parser.add_argument(
        "--critical-severities",
        nargs="+",
        default=critical_default,
        help="List of severities to treat as critical",
    )
    args = parser.parse_args()

    if not args.report_path:
        print("::error::Report path not specified via --report-path or REPORT_PATH_JSON", file=sys.stderr)
        sys.exit(1)

    try:
        process(
            report_path=args.report_path,
            module_name=args.module_name,
            summary_path=args.summary_path,
            output_path=args.output_path,
            critical_severities=args.critical_severities,
        )
    except Exception as e:
        print(f"::error::Failed to process lint report: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()

