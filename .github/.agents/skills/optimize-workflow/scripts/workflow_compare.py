#!/usr/bin/env python3
import argparse
import json
import os
import re
import sys

def parse_args():
    parser = argparse.ArgumentParser(description="Compare GitHub Actions Workflow Run Trials")
    parser.add_argument("trial_before", help="Name of the base trial (or workflow/trial path)")
    parser.add_argument("trial_after", help="Name of the new trial (or workflow/trial path)")
    return parser.parse_args()


def get_trials_dir():
    # Trial directories live directly inside the optimize-workflow skill directory.
    return os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))


def sanitize_dir_name(name):
    if not name:
        return "unknown"
    name = name.replace('/', '_').replace('\\', '_')
    while '..' in name:
        name = name.replace('..', '_')
    return name.strip()


def split_workflow_trial(trial_input):
    parts = trial_input.replace('\\', '/').split('/')
    if len(parts) == 2:
        return sanitize_dir_name(parts[0]), sanitize_dir_name(parts[1])
    return None, None


def find_report(trials_dir, workflow_name, trial_name):
    trial_dir = os.path.join(trials_dir, workflow_name, trial_name)
    report_path = os.path.join(trial_dir, "report.json")
    if os.path.isfile(report_path):
        return report_path
    legacy_path = os.path.join(trials_dir, workflow_name, f"{trial_name}.json")
    if os.path.isfile(legacy_path):
        return legacy_path
    return None


def search_trial(trials_dir, trial_name):
    safe_name = sanitize_dir_name(trial_name)
    candidates = []
    for entry in os.listdir(trials_dir):
        workflow_dir = os.path.join(trials_dir, entry)
        if not os.path.isdir(workflow_dir):
            continue
        report_path = find_report(trials_dir, entry, safe_name)
        if report_path:
            candidates.append((entry, report_path))
    if len(candidates) == 1:
        return candidates[0]
    if len(candidates) > 1:
        workflows = [c[0] for c in candidates]
        raise ValueError(f"Trial {trial_name!r} found in multiple workflows: {workflows}")
    raise FileNotFoundError(f"Could not find report for trial {trial_name!r} under {trials_dir}")


def resolve_trial(trials_dir, trial_input):
    workflow_name, trial_name = split_workflow_trial(trial_input)
    if workflow_name is not None:
        report_path = find_report(trials_dir, workflow_name, trial_name)
        if not report_path:
            raise FileNotFoundError(f"Could not find report for {trial_input!r}")
        return workflow_name, trial_name, report_path
    workflow_name, report_path = search_trial(trials_dir, trial_input)
    return workflow_name, sanitize_dir_name(trial_input), report_path

def load_report(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        return json.load(f)

def duration_to_seconds(dur_str):
    if not dur_str or dur_str == "N/A":
        return None
    try:
        parts = dur_str.split(':')
        if len(parts) == 3:
            return int(parts[0])*3600 + int(parts[1])*60 + int(parts[2])
        elif len(parts) == 2:
            return int(parts[0])*60 + int(parts[1])
        return int(parts[0])
    except Exception:
        return None

def seconds_to_duration(secs):
    sign = "-" if secs < 0 else "+" if secs > 0 else ""
    secs = abs(secs)
    h = secs // 3600
    m = (secs % 3600) // 60
    s = secs % 60
    return f"{sign}{h}:{m:02d}:{s:02d}"

def compare_durations(dur1, dur2):
    sec1 = duration_to_seconds(dur1)
    sec2 = duration_to_seconds(dur2)
    if sec1 is None or sec2 is None:
        return "N/A", "N/A"
    
    delta = sec2 - sec1
    diff_str = seconds_to_duration(delta)
    if delta == 0:
        diff_str = "0:00:00"
        
    if sec1 == 0:
        pct_str = "0.0%"
    else:
        pct = (delta / sec1) * 100
        pct_str = f"{pct:+.1f}%" if delta != 0 else "0.0%"
        
    return diff_str, pct_str

def parse_cost(cost_str):
    if not cost_str or cost_str == "N/A":
        return None
    try:
        return float(cost_str.replace('$', '').strip())
    except Exception:
        return None

def compare_costs(cost1, cost2):
    val1 = parse_cost(cost1)
    val2 = parse_cost(cost2)
    if val1 is None or val2 is None:
        return "N/A", "N/A"
        
    delta = val2 - val1
    if delta > 0:
        diff_str = f"+${delta:.4f}"
    elif delta < 0:
        diff_str = f"-${abs(delta):.4f}"
    else:
        diff_str = "$0.0000"
        
    if val1 == 0:
        pct_str = "0.0%"
    else:
        pct = (delta / val1) * 100
        pct_str = f"{pct:+.1f}%" if delta != 0 else "0.0%"
        
    return diff_str, pct_str

def parse_value_and_unit(avg_str):
    if not avg_str:
        return None, None
    match = re.match(r'^([\d\.]+)\s*(.*)$', avg_str.strip())
    if match:
        try:
            val = float(match.group(1))
            unit = match.group(2).strip()
            return val, unit
        except Exception:
            pass
    return None, None

def format_float(val):
    s = f"{val:.2f}"
    if s.endswith('.00'):
        return s[:-3]
    if s.endswith('0'):
        return s[:-1]
    return s

def compare_metrics(metric1, metric2):
    if not isinstance(metric1, dict) or not isinstance(metric2, dict):
        return "N/A"
        
    avg1 = metric1.get("avg")
    avg2 = metric2.get("avg")
    if not avg1 or not avg2:
        return "N/A"
        
    val1, unit1 = parse_value_and_unit(avg1)
    val2, unit2 = parse_value_and_unit(avg2)
    
    if val1 is not None and val2 is not None:
        delta = val2 - val1
        sign = "+" if delta > 0 else "-" if delta < 0 else ""
        unit = unit1 or unit2
        unit_suffix = f" {unit}" if unit else ""
        
        delta_str = format_float(delta)
        diff_val = f" ({sign}{abs(delta):.2f}{unit_suffix})" if delta != 0 else ""
        # clean up trailing zeros in diff_val too
        if diff_val:
            diff_val = diff_val.replace(f"{abs(delta):.2f}", format_float(abs(delta)))
            
        return f"avg: {avg1} -> {avg2}{diff_val}"
    
    return f"avg: {avg1} -> {avg2}"

def normalize_name(name):
    if not name:
        return ""
    name = name.lower()
    return re.sub(r'[^a-z0-9]', '', name)

def matches_job_name(name1, name2):
    a = normalize_name(name1)
    b = normalize_name(name2)
    if not a or not b:
        return False
    return (
        a.startswith(b)
        or b.startswith(a)
        or a.endswith(b)
        or b.endswith(a)
        or ((len(a) > 10 and len(b) > 10) and (a in b or b in a))
    )


def generate_comparison(data1, data2):
    run1 = data1.get("run", {})
    run2 = data2.get("run", {})
    
    run_dur_diff, run_dur_pct = compare_durations(run1.get("runtime"), run2.get("runtime"))
    
    lines = [
        "# Workflow Trial Comparison",
        f"- **Base Run ID (Trial 1)**: `{run1.get('id')}` (Status: `{run1.get('conclusion')}`, Runtime: `{run1.get('runtime')}`)",
        f"- **New Run ID (Trial 2)**: `{run2.get('id')}` (Status: `{run2.get('conclusion')}`, Runtime: `{run2.get('runtime')}`)",
        f"- **Runtime Delta**: `{run_dur_diff}` ({run_dur_pct})",
        "",
        "## Jobs Comparison"
    ]
    
    jobs1 = data1.get("jobs", [])
    jobs2 = data2.get("jobs", [])
    
    matched_jobs = []
    unmatched_jobs2 = list(jobs2)
    
    for job1 in jobs1:
        name1 = job1.get("name", "")
        # Find match in jobs2
        match_job2 = None
        # 1. Exact normalized match
        for job2 in jobs2:
            if normalize_name(name1) == normalize_name(job2.get("name", "")):
                match_job2 = job2
                break
        # 2. Fallback prefix match
        if not match_job2:
            for job2 in jobs2:
                if matches_job_name(name1, job2.get("name", "")):
                    match_job2 = job2
                    break
                    
        if match_job2:
            matched_jobs.append((job1, match_job2))
            if match_job2 in unmatched_jobs2:
                unmatched_jobs2.remove(match_job2)
        else:
            matched_jobs.append((job1, None))
            
    # Add unmatched jobs from Trial 2
    for job2 in unmatched_jobs2:
        matched_jobs.append((None, job2))
        
    for j1, j2 in matched_jobs:
        if j1 and j2:
            name = j1.get("name")
            dur_diff, dur_pct = compare_durations(j1.get("duration"), j2.get("duration"))
            
            m1 = j1.get("metrics", {})
            m2 = j2.get("metrics", {})
            
            cost_diff, cost_pct = compare_costs(m1.get("Cost"), m2.get("Cost"))
            
            lines.append(f"\n### Job: {name}")
            lines.append(f"- **Status**: `{j1.get('conclusion')}` -> `{j2.get('conclusion')}`")
            lines.append(f"- **Runner**: `{j1.get('runner', {}).get('labels', 'None')}` -> `{j2.get('runner', {}).get('labels', 'None')}`")
            lines.append(f"- **Runner Name**: `{j1.get('runner', {}).get('name', 'Unknown')}` -> `{j2.get('runner', {}).get('name', 'Unknown')}`")
            lines.append("")
            lines.append("| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |")
            lines.append("| --- | --- | --- | --- |")
            lines.append(f"| **Duration** | {j1.get('duration')} | {j2.get('duration')} | {dur_diff} ({dur_pct}) |")
            
            if m1.get("Instance Type") or m2.get("Instance Type"):
                lines.append(f"| **Instance Type** | {m1.get('Instance Type', 'N/A')} | {m2.get('Instance Type', 'N/A')} | {m1.get('Instance Type', 'N/A')} -> {m2.get('Instance Type', 'N/A')} |")
            if m1.get("Instance Lifecycle") or m2.get("Instance Lifecycle"):
                lines.append(f"| **Lifecycle** | {m1.get('Instance Lifecycle', 'N/A')} | {m2.get('Instance Lifecycle', 'N/A')} | - |")
            if m1.get("Cost") or m2.get("Cost"):
                lines.append(f"| **Cost** | {m1.get('Cost', 'N/A')} | {m2.get('Cost', 'N/A')} | {cost_diff} ({cost_pct}) |")
            if m1.get("Savings") or m2.get("Savings"):
                lines.append(f"| **Savings** | {m1.get('Savings', 'N/A')} | {m2.get('Savings', 'N/A')} | - |")
                
            cpu_diff = compare_metrics(m1.get("system.cpu.load_average.1m"), m2.get("system.cpu.load_average.1m"))
            if cpu_diff != "N/A":
                lines.append(f"| **CPU 1m (Avg)** | {m1.get('system.cpu.load_average.1m', {}).get('avg', 'N/A')} | {m2.get('system.cpu.load_average.1m', {}).get('avg', 'N/A')} | {cpu_diff} |")
                
            mem_diff = compare_metrics(m1.get("system.memory.utilization (used)"), m2.get("system.memory.utilization (used)"))
            if mem_diff != "N/A":
                lines.append(f"| **Memory (Avg)** | {m1.get('system.memory.utilization (used)', {}).get('avg', 'N/A')} | {m2.get('system.memory.utilization (used)', {}).get('avg', 'N/A')} | {mem_diff} |")
                
            disk_diff = compare_metrics(m1.get("disk_io_mb"), m2.get("disk_io_mb"))
            if disk_diff != "N/A":
                lines.append(f"| **Disk I/O (Avg)** | {m1.get('disk_io_mb', {}).get('avg', 'N/A')} | {m2.get('disk_io_mb', {}).get('avg', 'N/A')} | {disk_diff} |")
                
            net_diff = compare_metrics(m1.get("network_io_mb"), m2.get("network_io_mb"))
            if net_diff != "N/A":
                lines.append(f"| **Network I/O (Avg)** | {m1.get('network_io_mb', {}).get('avg', 'N/A')} | {m2.get('network_io_mb', {}).get('avg', 'N/A')} | {net_diff} |")
                
        elif j1:
            lines.append(f"\n### Job: {j1.get('name')} (Removed)")
            lines.append("- *This job was only present in Trial 1 (Base).*")
        elif j2:
            lines.append(f"\n### Job: {j2.get('name')} (Added)")
            lines.append("- *This job was only present in Trial 2 (New).*")
            
    return "\n".join(lines)

def main():
    args = parse_args()
    trials_dir = get_trials_dir()

    try:
        workflow1, trial1, path1 = resolve_trial(trials_dir, args.trial_before)
        workflow2, trial2, path2 = resolve_trial(trials_dir, args.trial_after)
    except Exception as e:
        print(f"Error locating trial reports: {e}", file=sys.stderr)
        sys.exit(1)

    if workflow1 != workflow2:
        print(f"Error: trials must belong to the same workflow (found {workflow1!r} and {workflow2!r})", file=sys.stderr)
        sys.exit(1)

    try:
        data1 = load_report(path1)
        data2 = load_report(path2)
    except Exception as e:
        print(f"Error loading JSON reports: {e}", file=sys.stderr)
        sys.exit(1)

    report = generate_comparison(data1, data2)
    print(report)

    out_dir = os.path.join(trials_dir, workflow1)
    os.makedirs(out_dir, exist_ok=True)
    out_file = os.path.join(out_dir, f"{trial1}-{trial2}-comparison.md")
    try:
        with open(out_file, 'w', encoding='utf-8') as f:
            f.write(report)
        print(f"Comparison report saved to: {out_file}", file=sys.stderr)
    except Exception as e:
        print(f"Error saving report to {out_file}: {e}", file=sys.stderr)

if __name__ == "__main__":
    main()
