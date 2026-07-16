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
    parser.add_argument("--repo", help="GitHub repository (owner/repo). Auto-detected if not specified.")
    parser.add_argument("--token", help="GitHub token. Defaults to GITHUB_TOKEN env var.")
    parser.add_argument("--poll-interval", type=int, default=10, help="Interval (seconds) to poll workflow run progress.")
    parser.add_argument("--format", choices=["markdown", "json"], default="markdown", help="Output format. Defaults to markdown.")
    parser.add_argument("--out-file", help="Path to write the output report.")
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
    name = name.lower()
    name = name.replace('/', '_')
    name = re.sub(r'\s+', '', name)
    name = name.replace('...', '')
    return name.rstrip('.')

def matches_job_name(api_name, log_job_name):
    a = normalize_name(api_name)
    b = normalize_name(log_job_name)
    return a.startswith(b) or b.startswith(a)

def generate_report(run_data, jobs, metrics, log_dir, format_type):
    start_time = run_data.get("run_started_at") or run_data.get("created_at")
    end_time = run_data.get("updated_at")
    total_runtime = format_duration(start_time, end_time)

    if format_type == "json":
        jobs_json = []
        for job in jobs:
            job_dur = format_duration(job.get("started_at"), job.get("completed_at"))
            labels_list = job.get("labels", [])
            runner_name = job.get("runner_name") or "Unknown"
            
            # Find metrics
            job_name = job.get('name', '')
            job_metrics = None
            for k, v in metrics.items():
                if normalize_name(job_name) == normalize_name(k):
                    job_metrics = v
                    break
            if not job_metrics:
                for k, v in metrics.items():
                    if matches_job_name(job_name, k):
                        job_metrics = v
                        break
                        
            job_entry = {
                "name": job.get("name"),
                "status": job.get("status"),
                "conclusion": job.get("conclusion"),
                "runner": {
                    "labels": labels_list,
                    "name": runner_name
                },
                "started_at": job.get("started_at"),
                "completed_at": job.get("completed_at"),
                "duration": job_dur,
                "metrics": job_metrics or {}
            }
            jobs_json.append(job_entry)
            
        report_data = {
            "run": {
                "id": run_data.get("id"),
                "status": run_data.get("status"),
                "conclusion": run_data.get("conclusion"),
                "runtime": total_runtime
            },
            "logs_dir": log_dir,
            "jobs": jobs_json
        }
        return json.dumps(report_data, indent=2)
        
    else:  # markdown
        lines = [
            "# Workflow Run Summary",
            f"- **ID**: `{run_data.get('id')}`",
            f"- **Status**: `{run_data.get('status')}`",
            f"- **Conclusion**: `{run_data.get('conclusion')}`",
            f"- **Runtime**: `{total_runtime}`",
            f"- **Logs Directory**: `{log_dir}`",
            "",
            "## Jobs Summary"
        ]
        for job in jobs:
            job_dur = format_duration(job.get("started_at"), job.get("completed_at"))
            labels_list = job.get("labels", [])
            labels_str = ", ".join(labels_list) if labels_list else "None"
            runner_name = job.get("runner_name") or "Unknown"
            
            lines.append(f"\n### Job: {job.get('name')}")
            lines.append(f"- **Status**: `{job.get('status')}`")
            lines.append(f"- **Conclusion**: `{job.get('conclusion')}`")
            lines.append(f"- **Runner**: `{labels_str}` ({runner_name})")
            lines.append(f"- **Start**: `{job.get('started_at')}`")
            lines.append(f"- **End**: `{job.get('completed_at')}`")
            lines.append(f"- **Duration**: `{job_dur}`")
            
            # Look up job metrics
            job_name = job.get('name', '')
            job_metrics = None
            for k, v in metrics.items():
                if normalize_name(job_name) == normalize_name(k):
                    job_metrics = v
                    break
            if not job_metrics:
                for k, v in metrics.items():
                    if matches_job_name(job_name, k):
                        job_metrics = v
                        break
                        
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
    
    log_stderr(f"Monitoring workflow run {args.run_id} in {repo}...")
    
    try:
        run_data = wait_for_run(repo, args.run_id, token)
    except RuntimeError as e:
        log_stderr(f"Error: {e}")
        sys.exit(1)
        
    log_stderr("Workflow run found. Waiting for completion...")
    run_data = wait_for_completion(repo, args.run_id, token, args.poll_interval)
    
    # Download logs to persistent temp dir
    import tempfile
    log_dir = tempfile.mkdtemp(prefix=f"workflow-logs-{args.run_id}-")
    log_stderr(f"Downloading logs to: {log_dir}...")
    
    metrics = {}
    downloaded = download_logs(repo, args.run_id, token, log_dir)
    if downloaded:
        metrics = parse_job_logs(log_dir)
            
    jobs = fetch_jobs(repo, args.run_id, token)
    
    # Generate and print report
    report = generate_report(run_data, jobs, metrics, log_dir, args.format)
    
    # Write to file if specified
    if args.out_file:
        try:
            with open(args.out_file, 'w', encoding='utf-8') as f:
                f.write(report)
            log_stderr(f"Report successfully saved to: {args.out_file}")
            if args.format == 'json':
                log_stderr(f"Try exploring with: jq . {args.out_file}")
        except Exception as e:
            log_stderr(f"Error saving report to {args.out_file}: {e}")
    else:
        print(report)

if __name__ == "__main__":
    main()
