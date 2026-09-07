#!/usr/bin/env python3
import argparse
import datetime
import json
import os
import re
import sys
import time
import urllib.request
import urllib.error
import zipfile
import io

def parse_args():
    parser = argparse.ArgumentParser(description="Monitor GitHub Actions Workflow Run")
    parser.add_argument("run_id", type=int, help="GitHub Actions Workflow Run ID")
    parser.add_argument("trial_name", help="Trial name used to create the output directory.")
    parser.add_argument("--repo", help="GitHub repository (owner/repo). Auto-detected if not specified.")
    parser.add_argument("--token", help="GitHub token. Defaults to GITHUB_TOKEN env var.")
    parser.add_argument("--poll-interval", type=int, default=10, help="Interval (seconds) to poll workflow run progress.")
    return parser.parse_args()

def log_stderr(msg):
    timestamp = datetime.datetime.now().strftime('%H:%M:%S')
    print(f"[{timestamp}] {msg}", file=sys.stderr)

def detect_repo():
    try:
        import subprocess
        res = subprocess.run(['git', 'remote', 'get-url', 'origin'], capture_output=True, text=True, check=True)
        url = res.stdout.strip()
        match = re.search(r'github\.com[:/]([^/]+/[^/]+?)(?:\.git)?$', url)
        if match:
            return match.group(1)
    except Exception:
        pass
    return "smartcontractkit/chainlink"


def get_trials_dir():
    # Trial directories live directly inside the optimize-workflow skill directory.
    return os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'trials'))


def sanitize_dir_name(name):
    if not name:
        return "unknown"
    name = name.replace('/', '_').replace('\\', '_')
    while '..' in name:
        name = name.replace('..', '_')
    return name.strip()


def derive_workflow_name(run_data):
    path = run_data.get('path')
    if path:
        name = os.path.basename(path)
    else:
        name = run_data.get('name') or 'unknown'
    return sanitize_dir_name(name)


def get_headers(token):
    headers = {
        'Accept': 'application/vnd.github+json',
        'X-GitHub-Api-Version': '2022-11-28'
    }
    if token:
        headers['Authorization'] = f'Bearer {token}'
    return headers

def fetch_json(url, token):
    req = urllib.request.Request(url, headers=get_headers(token))
    with urllib.request.urlopen(req) as response:
        return json.loads(response.read().decode('utf-8'))

def wait_for_run(repo, run_id, token, timeout=10, poll_interval=2):
    url = f"https://api.github.com/repos/{repo}/actions/runs/{run_id}"
    start_time = time.time()
    while time.time() - start_time < timeout:
        try:
            return fetch_json(url, token)
        except urllib.error.HTTPError as e:
            if e.code == 404:
                log_stderr(f"Workflow run {run_id} not found yet. Retrying...")
                time.sleep(poll_interval)
            else:
                raise e
    raise RuntimeError(f"Workflow run {run_id} not found within {timeout} seconds.")

def wait_for_completion(repo, run_id, token, poll_interval=5):
    url = f"https://api.github.com/repos/{repo}/actions/runs/{run_id}"
    last_status = None
    while True:
        run_data = fetch_json(url, token)
        status = run_data.get("status")
        if status == "completed" or run_data.get("conclusion") is not None:
            return run_data
        if status != last_status:
            log_stderr(f"Workflow run status: {status}. Waiting...")
            last_status = status
        time.sleep(poll_interval)

def fetch_jobs(repo, run_id, token):
    jobs = []
    page = 1
    while True:
        url = f"https://api.github.com/repos/{repo}/actions/runs/{run_id}/jobs?per_page=100&page={page}"
        try:
            data = fetch_json(url, token)
        except Exception as e:
            print(f"Warning: Could not fetch jobs page {page}: {e}", file=sys.stderr)
            break
        page_jobs = data.get("jobs", [])
        if not page_jobs:
            break
        jobs.extend(page_jobs)
        if len(page_jobs) < 100:
            break
        page += 1
    return jobs

def download_logs(repo, run_id, token, dest_dir):
    url = f"https://api.github.com/repos/{repo}/actions/runs/{run_id}/logs"
    req = urllib.request.Request(url, headers=get_headers(token))
    try:
        with urllib.request.urlopen(req) as response:
            zip_data = response.read()
            with zipfile.ZipFile(io.BytesIO(zip_data)) as z:
                z.extractall(dest_dir)
            return True
    except Exception as e:
        print(f"Warning: Could not download logs: {e}", file=sys.stderr)
        return False

def parse_job_logs(log_dir):
    results = {}
    log_files = []
    for root, _, files in os.walk(log_dir):
        for file in files:
            if file.endswith('.txt'):
                log_files.append(os.path.join(root, file))
                
    for filepath in log_files:
        # Extract job name from filename (e.g., 41_Job Name _ subjob.txt -> Job Name / subjob)
        filename = os.path.basename(filepath)
        if filename.endswith('.txt'):
            filename = filename[:-4]
        match = re.match(r'^\d+_(.*)$', filename)
        if match:
            filename = match.group(1)
        job_name = filename.replace(' _ ', ' / ')
            
        try:
            with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
        except Exception:
            continue
            
        # Only parse if file contains runson indicators
        if "## Execution Cost Summary" not in content and "system.cpu.load_average.1m" not in content:
            continue
            
        if job_name not in results:
            results[job_name] = {}
            
        metrics = results[job_name]
        
        # 1. Parse cost summary table
        table_match = re.search(r'## Execution Cost Summary\s*([\s\S]*?)(?:\n\n|\Z)', content)
        if table_match:
            table_text = table_match.group(1)
            for line in table_text.splitlines():
                if '|' in line:
                    parts = [p.strip() for p in line.split('|')]
                    if len(parts) >= 3:
                        key = parts[1]
                        val = parts[2]
                        if key in [
                            "Instance Type", "Instance Lifecycle", "Region", "Platform",
                            "Arch", "Az", "Zone ID", "Duration", "Cost",
                            "GitHub equivalent cost", "Savings"
                        ]:
                            metrics[key] = val
                            
        # 2. Parse system metrics
        metric_patterns = [
            (r'system\.cpu\.load_average\.1m', 'system.cpu.load_average.1m'),
            (r'system\.cpu\.load_average\.5m', 'system.cpu.load_average.5m'),
            (r'system\.memory\.utilization\s*\(used\)', 'system.memory.utilization (used)')
        ]
        for pattern_regex, key_name in metric_patterns:
            regex = pattern_regex + r'\s+min:\s*(\S+)\s+max:\s*(\S+)\s+avg:\s*([\S\s]+?)(?=\n|\Z)'
            match = re.search(regex, content)
            if match:
                min_val = match.group(1).strip()
                max_val = match.group(2).strip()
                avg_val = match.group(3).strip()
                metrics[key_name] = {"min": min_val, "max": max_val, "avg": avg_val}
                
        # 3. Parse Disk/Network I/O
        disk_io_match = re.search(r'Disk I/O by Direction\s*\(MB\)\s*\(min:\s*(\S+),\s*max:\s*(\S+),\s*avg:\s*([\S\s]+?)\)', content)
        if disk_io_match:
            metrics['disk_io_mb'] = {"min": disk_io_match.group(1).strip(), "max": disk_io_match.group(2).strip(), "avg": disk_io_match.group(3).strip()}
            
        disk_ops_match = re.search(r'Disk Operations by Direction\s*\(min:\s*(\S+),\s*max:\s*(\S+),\s*avg:\s*([\S\s]+?)\)', content)
        if disk_ops_match:
            metrics['disk_ops'] = {"min": disk_ops_match.group(1).strip(), "max": disk_ops_match.group(2).strip(), "avg": disk_ops_match.group(3).strip()}
            
        net_io_match = re.search(r'Network I/O by Direction\s*\(MB\)\s*\(min:\s*(\S+),\s*max:\s*(\S+),\s*avg:\s*([\S\s]+?)\)', content)
        if net_io_match:
            metrics['network_io_mb'] = {"min": net_io_match.group(1).strip(), "max": net_io_match.group(2).strip(), "avg": net_io_match.group(3).strip()}
                
    # Clean up empty jobs
    results = {k: v for k, v in results.items() if v}
    return results

def parse_duration_seconds(start_str, end_str):
    if not start_str or not end_str:
        return 0
    try:
        start = datetime.datetime.strptime(start_str.replace("Z", ""), "%Y-%m-%dT%H:%M:%S")
        end = datetime.datetime.strptime(end_str.replace("Z", ""), "%Y-%m-%dT%H:%M:%S")
        delta = end - start
        return int(delta.total_seconds())
    except Exception:
        return 0

def format_duration(start_str, end_str):
    if not start_str or not end_str:
        return "N/A"
    try:
        # Handle GitHub ISO timestamps, e.g. "2026-07-16T14:09:38Z"
        # Standard lib handles 'Z' suffix using replace/strptime
        fmt = "%Y-%m-%dT%H:%M:%SZ"
        start = datetime.datetime.strptime(start_str.replace("Z", ""), "%Y-%m-%dT%H:%M:%S")
        end = datetime.datetime.strptime(end_str.replace("Z", ""), "%Y-%m-%dT%H:%M:%S")
        delta = end - start
        return str(delta)
    except Exception:
        return "N/A"

def normalize_name(name):
    if not name:
        return ""
    name = name.lower()
    return re.sub(r'[^a-z0-9]', '', name)

def matches_job_name(api_name, log_job_name):
    a = normalize_name(api_name)
    b = normalize_name(log_job_name)
    if not a or not b:
        return False
    return (
        a.startswith(b)
        or b.startswith(a)
        or a.endswith(b)
        or b.endswith(a)
        or ((len(a) > 10 and len(b) > 10) and (a in b or b in a))
    )

def find_job_metrics(job_name, metrics):
    for k, v in metrics.items():
        if normalize_name(job_name) == normalize_name(k):
            return v
    for k, v in metrics.items():
        if matches_job_name(job_name, k):
            return v
    return None

def get_job_cost(metrics):
    return metrics.get("Cost", "N/A") if metrics else "N/A"

def get_run_url(run_data):
    if not run_data:
        return None
    if run_data.get("html_url"):
        return run_data.get("html_url")
    run_id = run_data.get("id")
    if run_id:
        return f"https://github.com/smartcontractkit/chainlink/actions/runs/{run_id}"
    return None

def get_job_url(job, run_url=None):
    if not job:
        return None
    if job.get("html_url"):
        return job.get("html_url")
    job_id = job.get("id")
    if job_id and run_url:
        return f"{run_url}/job/{job_id}"
    return None

def generate_report(run_data, jobs, metrics, log_dir, format_type):
    start_time = run_data.get("run_started_at") or run_data.get("created_at")
    end_time = run_data.get("updated_at")
    total_runtime = format_duration(start_time, end_time)
    total_runtime_sec = parse_duration_seconds(start_time, end_time)
    run_url = get_run_url(run_data)

    # Compute workflow total cost
    total_cost_val = 0.0
    has_cost_data = False
    for job in jobs:
        jm = find_job_metrics(job.get('name', ''), metrics)
        c_str = get_job_cost(jm)
        if c_str != "N/A":
            try:
                val = float(c_str.replace('$', '').strip())
                total_cost_val += val
                has_cost_data = True
            except ValueError:
                pass

    total_workflow_cost = f"${total_cost_val:.4f}" if has_cost_data else "N/A"

    if format_type == "json":
        jobs_json = []
        slowest = []
        durations = []
        for job in jobs:
            job_dur_sec = parse_duration_seconds(job.get("started_at"), job.get("completed_at"))
            if job_dur_sec > 0:
                durations.append(job_dur_sec)

        avg_dur_sec = int(sum(durations) / len(durations)) if durations else 0

        for job in jobs:
            job_dur = format_duration(job.get("started_at"), job.get("completed_at"))
            job_dur_sec = parse_duration_seconds(job.get("started_at"), job.get("completed_at"))
            labels_list = job.get("labels", [])
            runner_name = job.get("runner_name") or "Unknown"
            is_outlier = job_dur_sec > (1.5 * avg_dur_sec) if avg_dur_sec > 0 else False
            
            job_name = job.get('name', '')
            job_metrics = find_job_metrics(job_name, metrics)
            job_cost = get_job_cost(job_metrics)
            job_url = get_job_url(job, run_url)
                        
            job_entry = {
                "name": job.get("name"),
                "status": job.get("status"),
                "conclusion": job.get("conclusion"),
                "html_url": job_url,
                "runner": {
                    "labels": labels_list,
                    "name": runner_name
                },
                "started_at": job.get("started_at"),
                "completed_at": job.get("completed_at"),
                "duration": job_dur,
                "duration_seconds": job_dur_sec,
                "cost": job_cost,
                "is_outlier": is_outlier,
                "metrics": job_metrics or {}
            }
            jobs_json.append(job_entry)
            slowest.append({
                "name": job.get("name"),
                "duration": job_dur,
                "duration_seconds": job_dur_sec,
                "cost": job_cost,
                "html_url": job_url,
                "is_outlier": is_outlier,
                "status": job.get("status"),
                "conclusion": job.get("conclusion")
            })

        slowest.sort(key=lambda x: x["duration_seconds"], reverse=True)

        report_data = {
            "run": {
                "id": run_data.get("id"),
                "status": run_data.get("status"),
                "conclusion": run_data.get("conclusion"),
                "html_url": run_url,
                "runtime": total_runtime,
                "runtime_seconds": total_runtime_sec,
                "total_cost": total_workflow_cost,
                "avg_job_duration_seconds": avg_dur_sec
            },
            "logs_dir": log_dir,
            "slowest_jobs": slowest[:10],
            "jobs": jobs_json
        }
        return json.dumps(report_data, indent=2)
        
    else:  # markdown
        run_id_str = f"[{run_data.get('id')}]({run_url})" if run_url else f"`{run_data.get('id')}`"
        run_status_str = f"[{run_data.get('status')}]({run_url})" if run_url else f"`{run_data.get('status')}`"

        lines = [
            "# Workflow Run Summary",
            f"- **ID**: {run_id_str}",
            f"- **Status**: {run_status_str}",
            f"- **Conclusion**: `{run_data.get('conclusion')}`",
            f"- **Runtime**: `{total_runtime}`",
            f"- **Total Cost**: `{total_workflow_cost}`",
            f"- **Logs Directory**: `{log_dir}`",
            "",
            "## Longest Jobs (Bottlenecks)"
        ]

        # Calculate durations and sort jobs descending
        sorted_jobs = []
        for j in jobs:
            dur_sec = parse_duration_seconds(j.get("started_at"), j.get("completed_at"))
            sorted_jobs.append((dur_sec, j))
        sorted_jobs.sort(key=lambda x: x[0], reverse=True)

        lines.append("| Job Name | Duration | Cost | Status | Conclusion | Runner |")
        lines.append("| --- | --- | --- | --- | --- | --- |")
        for dur_sec, j in sorted_jobs[:10]:
            dur_str = format_duration(j.get("started_at"), j.get("completed_at"))
            jm = find_job_metrics(j.get('name', ''), metrics)
            c_str = get_job_cost(jm)
            labels_list = j.get("labels", [])
            labels_str = ", ".join(labels_list) if labels_list else "None"
            job_url = get_job_url(j, run_url)
            job_name_str = f"[{j.get('name')}]({job_url})" if job_url else j.get('name')
            job_status_str = f"[{j.get('status')}]({job_url})" if job_url else f"`{j.get('status')}`"
            lines.append(f"| {job_name_str} | `{dur_str}` | `{c_str}` | {job_status_str} | `{j.get('conclusion')}` | `{labels_str}` |")

        lines.append("\n## Jobs Summary")

        for job in jobs:
            job_dur = format_duration(job.get("started_at"), job.get("completed_at"))
            job_name = job.get('name', '')
            job_metrics = find_job_metrics(job_name, metrics)
            job_cost = get_job_cost(job_metrics)
            job_url = get_job_url(job, run_url)
            job_name_str = f"[{job.get('name')}]({job_url})" if job_url else job.get('name')
            job_status_str = f"[{job.get('status')}]({job_url})" if job_url else f"`{job.get('status')}`"
            labels_list = job.get("labels", [])
            labels_str = ", ".join(labels_list) if labels_list else "None"
            runner_name = job.get("runner_name") or "Unknown"
            
            lines.append(f"\n### Job: {job_name_str}")
            lines.append(f"- **Status**: {job_status_str}")
            lines.append(f"- **Conclusion**: `{job.get('conclusion')}`")
            lines.append(f"- **Runner**: `{labels_str}` ({runner_name})")
            lines.append(f"- **Start**: `{job.get('started_at')}`")
            lines.append(f"- **End**: `{job.get('completed_at')}`")
            lines.append(f"- **Duration**: `{job_dur}`")
            lines.append(f"- **Cost**: `{job_cost}`")
            
            if job_metrics:
                lines.append("- **Runner Details**:")
                lines.append("  | metric | value |")
                lines.append("  | --- | --- |")
                for key in [
                    "Instance Type", "Instance Lifecycle", "Region", "Platform",
                    "Arch", "Az", "Zone ID", "Duration", "Cost",
                    "GitHub equivalent cost", "Savings"
                ]:
                    if key in job_metrics:
                        lines.append(f"  | {key} | {job_metrics[key]} |")
                        
                lines.append("- **CPU**:")
                if 'system.cpu.load_average.1m' in job_metrics:
                    m = job_metrics['system.cpu.load_average.1m']
                    lines.append(f"  - `system.cpu.load_average.1m`: min: {m['min']}, max: {m['max']}, avg: {m['avg']}")
                if 'system.cpu.load_average.5m' in job_metrics:
                    m = job_metrics['system.cpu.load_average.5m']
                    lines.append(f"  - `system.cpu.load_average.5m`: min: {m['min']}, max: {m['max']}, avg: {m['avg']}")
                    
                lines.append("- **Memory**:")
                if 'system.memory.utilization (used)' in job_metrics:
                    m = job_metrics['system.memory.utilization (used)']
                    lines.append(f"  - `system.memory.utilization (used)`: min: {m['min']}, max: {m['max']}, avg: {m['avg']}")
                    
                lines.append("- **Disk**:")
                if 'disk_io_mb' in job_metrics:
                    m = job_metrics['disk_io_mb']
                    lines.append(f"  - `Disk I/O by Direction`: min: {m['min']}, max: {m['max']}, avg: {m['avg']}")
                if 'disk_ops' in job_metrics:
                    m = job_metrics['disk_ops']
                    lines.append(f"  - `Disk Operations by Direction`: min: {m['min']}, max: {m['max']}, avg: {m['avg']}")
                    
                lines.append("- **I/O**:")
                if 'network_io_mb' in job_metrics:
                    m = job_metrics['network_io_mb']
                    lines.append(f"  - `Network I/O by Direction`: min: {m['min']}, max: {m['max']}, avg: {m['avg']}")
                    
        return "\n".join(lines)

def main():
    args = parse_args()
    token = args.token or os.environ.get("GITHUB_TOKEN")
    repo = args.repo or detect_repo()
    trial_name = sanitize_dir_name(args.trial_name)

    log_stderr(f"Monitoring workflow run {args.run_id} in {repo}...")

    try:
        run_data = wait_for_run(repo, args.run_id, token)
    except RuntimeError as e:
        log_stderr(f"Error: {e}")
        sys.exit(1)

    workflow_name = derive_workflow_name(run_data)
    trials_dir = get_trials_dir()
    output_dir = os.path.join(trials_dir, workflow_name, trial_name)
    logs_dir = os.path.join(output_dir, "logs")

    try:
        os.makedirs(output_dir, exist_ok=True)
        os.makedirs(logs_dir, exist_ok=True)
    except Exception as e:
        log_stderr(f"Error creating output directories: {e}")
        sys.exit(1)

    log_stderr("Workflow run found. Waiting for completion...")
    run_data = wait_for_completion(repo, args.run_id, token, args.poll_interval)

    log_stderr(f"Downloading logs to: {logs_dir}...")
    metrics = {}
    downloaded = download_logs(repo, args.run_id, token, logs_dir)
    if downloaded:
        metrics = parse_job_logs(logs_dir)

    jobs = fetch_jobs(repo, args.run_id, token)

    json_path = os.path.join(output_dir, "report.json")
    md_path = os.path.join(output_dir, "report.md")

    try:
        json_report = generate_report(run_data, jobs, metrics, logs_dir, "json")
        with open(json_path, 'w', encoding='utf-8') as f:
            f.write(json_report)
        log_stderr(f"JSON report saved to: {json_path}")
    except Exception as e:
        log_stderr(f"Error saving JSON report: {e}")

    try:
        md_report = generate_report(run_data, jobs, metrics, logs_dir, "markdown")
        with open(md_path, 'w', encoding='utf-8') as f:
            f.write(md_report)
        log_stderr(f"Markdown report saved to: {md_path}")
    except Exception as e:
        log_stderr(f"Error saving Markdown report: {e}")

    print(f"workflow: {workflow_name}")
    print(f"trial: {trial_name}")
    print(f"json: {json_path}")
    print(f"markdown: {md_path}")
    print(f"logs: {logs_dir}")

if __name__ == "__main__":
    main()
