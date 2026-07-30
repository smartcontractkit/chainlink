import unittest
from unittest.mock import patch, MagicMock
import io
import os
import tempfile
import zipfile
import json
import urllib.error
import urllib.request
import sys

# Import the stub module
import workflow_monitor

class TestWorkflowMonitor(unittest.TestCase):

    @patch('subprocess.run')
    def test_detect_repo_git(self, mock_run):
        mock_response = MagicMock()
        mock_run.return_value = mock_response

        # SSH format
        mock_response.stdout = "git@github.com:smartcontractkit/chainlink.git\n"
        self.assertEqual(workflow_monitor.detect_repo(), "smartcontractkit/chainlink")

        # HTTPS format
        mock_response.stdout = "https://github.com/smartcontractkit/chainlink.git\n"
        self.assertEqual(workflow_monitor.detect_repo(), "smartcontractkit/chainlink")

        # HTTPS without .git
        mock_response.stdout = "https://github.com/smartcontractkit/chainlink\n"
        self.assertEqual(workflow_monitor.detect_repo(), "smartcontractkit/chainlink")

        # Malicious subdomain containing github.com
        mock_response.stdout = "https://github.com.attacker.com/smartcontractkit/chainlink.git\n"
        self.assertEqual(workflow_monitor.detect_repo(), "smartcontractkit/chainlink")

        # Fallback on subprocess error
        mock_run.side_effect = Exception("git command failed")
        self.assertEqual(workflow_monitor.detect_repo(), "smartcontractkit/chainlink")

    @patch('urllib.request.urlopen')
    def test_fetch_json_success(self, mock_urlopen):
        mock_response = MagicMock()
        mock_response.read.return_value = b'{"status": "completed", "conclusion": "success"}'
        mock_urlopen.return_value.__enter__.return_value = mock_response
        
        result = workflow_monitor.fetch_json('https://api.github.com/dummy', 'fake_token')
        self.assertEqual(result['status'], 'completed')
        self.assertEqual(result['conclusion'], 'success')
        
    @patch('urllib.request.urlopen')
    def test_wait_for_run_found_immediately(self, mock_urlopen):
        mock_response = MagicMock()
        mock_response.read.return_value = b'{"id": 12345}'
        mock_urlopen.return_value.__enter__.return_value = mock_response
        
        result = workflow_monitor.wait_for_run('owner/repo', 12345, 'fake_token', timeout=5, poll_interval=1)
        self.assertEqual(result['id'], 12345)

    @patch('urllib.request.urlopen')
    @patch('time.sleep')
    @patch('builtins.print')
    def test_wait_for_run_not_found_then_found(self, mock_print, mock_sleep, mock_urlopen):
        # First call raises HTTPError 404, second call succeeds
        mock_response = MagicMock()
        mock_response.read.return_value = b'{"id": 12345}'
        
        req = urllib.request.Request('https://api.github.com/repos/owner/repo/actions/runs/12345')
        err = urllib.error.HTTPError(req.full_url, 404, 'Not Found', {}, None)
        
        # We need the context manager for the second call to succeed
        mock_cm = MagicMock()
        mock_cm.__enter__.return_value = mock_response
        
        mock_urlopen.side_effect = [err, mock_cm]
        
        result = workflow_monitor.wait_for_run('owner/repo', 12345, 'fake_token', timeout=5, poll_interval=1)
        self.assertEqual(result['id'], 12345)
        self.assertEqual(mock_urlopen.call_count, 2)
        mock_sleep.assert_called_once_with(1)
        self.assertEqual(mock_print.call_count, 1)
        args, kwargs = mock_print.call_args
        self.assertEqual(kwargs.get('file'), sys.stderr)
        import re
        self.assertTrue(re.search(r'^\[\d{2}:\d{2}:\d{2}\] Workflow run 12345 not found yet\. Retrying\.\.\.', args[0]))

    @patch('urllib.request.urlopen')
    @patch('time.sleep')
    @patch('time.time')
    @patch('builtins.print')
    def test_wait_for_run_timeout(self, mock_print, mock_time, mock_sleep, mock_urlopen):
        mock_time.side_effect = [100.0, 100.5, 101.5, 103.0]
        req = urllib.request.Request('https://api.github.com/repos/owner/repo/actions/runs/12345')
        err = urllib.error.HTTPError(req.full_url, 404, 'Not Found', {}, None)
        mock_urlopen.side_effect = err
        
        with self.assertRaises(RuntimeError):
            workflow_monitor.wait_for_run('owner/repo', 12345, 'fake_token', timeout=2, poll_interval=1)

    @patch('urllib.request.urlopen')
    @patch('time.sleep')
    @patch('builtins.print')
    def test_wait_for_completion(self, mock_print, mock_sleep, mock_urlopen):
        # Three checks: in_progress, in_progress, completed
        mock_resp_1 = MagicMock()
        mock_resp_1.read.return_value = b'{"status": "in_progress", "conclusion": null}'
        mock_cm_1 = MagicMock()
        mock_cm_1.__enter__.return_value = mock_resp_1
        
        mock_resp_2 = MagicMock()
        mock_resp_2.read.return_value = b'{"status": "in_progress", "conclusion": null}'
        mock_cm_2 = MagicMock()
        mock_cm_2.__enter__.return_value = mock_resp_2
        
        mock_resp_3 = MagicMock()
        mock_resp_3.read.return_value = b'{"status": "completed", "conclusion": "success"}'
        mock_cm_3 = MagicMock()
        mock_cm_3.__enter__.return_value = mock_resp_3
        
        mock_urlopen.side_effect = [mock_cm_1, mock_cm_2, mock_cm_3]
        
        result = workflow_monitor.wait_for_completion('owner/repo', 12345, 'fake_token', poll_interval=1)
        self.assertEqual(result['conclusion'], 'success')
        self.assertEqual(mock_urlopen.call_count, 3)
        self.assertEqual(mock_sleep.call_count, 2)
        
        # Only prints once for in_progress status
        self.assertEqual(mock_print.call_count, 1)
        args, kwargs = mock_print.call_args
        self.assertEqual(kwargs.get('file'), sys.stderr)
        
        import re
        self.assertTrue(re.search(r'^\[\d{2}:\d{2}:\d{2}\] Workflow run status: in_progress. Waiting...', args[0]))

    def test_parse_job_logs_extracted(self):
        # Create a temp directory structure mimicking extracted logs
        with tempfile.TemporaryDirectory() as tmp_dir:
            log_content = """
2026-07-16T17:16:03.3763159Z Mapped zone name us-east-2c to zone ID use2-az3
2026-07-16T17:16:03.5451437Z ## Execution Cost Summary
2026-07-16T17:16:03.5451730Z 
2026-07-16T17:16:03.5451928Z | metric                 | value           |
2026-07-16T17:16:03.5452417Z | ---------------------- | --------------- |
2026-07-16T17:16:03.5452778Z | Instance Type          | c6in.4xlarge    |
2026-07-16T17:16:03.5453143Z | Instance Lifecycle     | spot            |
2026-07-16T17:16:03.5453493Z | Region                 | us-east-2       |
2026-07-16T17:16:03.5453851Z | Platform               | Linux/UNIX      |
2026-07-16T17:16:03.5454200Z | Arch                   | x64             |
2026-07-16T17:16:03.5454536Z | Az                     | us-east-2c      |
2026-07-16T17:16:03.5454880Z | Zone ID                | use2-az3        |
2026-07-16T17:16:03.5455225Z | Duration               | 3.52 minutes    |
2026-07-16T17:16:03.5455559Z | Cost                   | $0.0179         |
2026-07-16T17:16:03.5456195Z | GitHub equivalent cost | $0.1680         |
2026-07-16T17:16:03.5456559Z | Savings                | $0.1501 (89.3%) |
2026-07-16T17:16:03.5456774Z 
2026-07-16T17:16:04.0020706Z   system.cpu.load_average.1m                     min: 0.06  max: 19.83  avg: 7.51 
2026-07-16T17:16:04.0021347Z   system.cpu.load_average.5m                     min: 0.02  max: 5.17  avg: 2.62 
2026-07-16T17:16:04.0021902Z   system.memory.utilization (used)               min: 1.53  max: 11.97  avg: 6.78 %
2026-07-16T17:16:04.0008487Z        Disk I/O by Direction (MB) (min: 44.63, max: 590.74, avg: 467.01 MB)
2026-07-16T17:16:04.0013105Z         Disk Operations by Direction (min: 757, max: 15175, avg: 8014.92 ops)
2026-07-16T17:16:04.0018858Z        Network I/O by Direction (MB) (min: 0.24, max: 5837.67, avg: 1220.81 MB)
"""
            # Create a flat file at root mimicking 15_[Job Name].txt
            with open(os.path.join(tmp_dir, "15_build_job.txt"), "w") as f:
                f.write(log_content)
                
            metrics = workflow_monitor.parse_job_logs(tmp_dir)
            
            self.assertIn("build_job", metrics)
            job_metrics = metrics["build_job"]
            self.assertEqual(job_metrics["Instance Type"], "c6in.4xlarge")
            self.assertEqual(job_metrics["Instance Lifecycle"], "spot")
            self.assertEqual(job_metrics["Savings"], "$0.1501 (89.3%)")
            
            # Check system metrics
            self.assertEqual(job_metrics["system.cpu.load_average.1m"], {"min": "0.06", "max": "19.83", "avg": "7.51"})
            self.assertEqual(job_metrics["system.cpu.load_average.5m"], {"min": "0.02", "max": "5.17", "avg": "2.62"})
            self.assertEqual(job_metrics["system.memory.utilization (used)"], {"min": "1.53", "max": "11.97", "avg": "6.78 %"})
            
            # Check disk/network metrics
            self.assertEqual(job_metrics["disk_io_mb"], {"min": "44.63", "max": "590.74", "avg": "467.01 MB"})
            self.assertEqual(job_metrics["disk_ops"], {"min": "757", "max": "15175", "avg": "8014.92 ops"})
            self.assertEqual(job_metrics["network_io_mb"], {"min": "0.24", "max": "5837.67", "avg": "1220.81 MB"})

    def test_generate_report_markdown(self):
        run_data = {"id": 12345, "status": "completed", "conclusion": "success", "run_started_at": "2026-07-16T17:16:52Z", "updated_at": "2026-07-16T17:21:52Z"}
        jobs = [{"name": "build", "status": "completed", "conclusion": "success", "labels": ["runs-on"], "runner_name": "runner-1", "started_at": "2026-07-16T17:16:52Z", "completed_at": "2026-07-16T17:18:52Z"}]
        metrics = {"build": {"Instance Type": "c6in.4xlarge", "disk_io_mb": {"min": "44", "max": "590", "avg": "467 MB"}}}
        
        report = workflow_monitor.generate_report(run_data, jobs, metrics, "/tmp/logs-dir", "markdown")
        self.assertIn("# Workflow Run Summary", report)
        self.assertIn("- **Logs Directory**: `/tmp/logs-dir`", report)
        self.assertIn("c6in.4xlarge", report)
        self.assertIn("Disk I/O by Direction", report)

    def test_generate_report_json(self):
        run_data = {"id": 12345, "status": "completed", "conclusion": "success", "run_started_at": "2026-07-16T17:16:52Z", "updated_at": "2026-07-16T17:21:52Z"}
        jobs = [{"name": "build", "status": "completed", "conclusion": "success", "labels": ["runs-on"], "runner_name": "runner-1", "started_at": "2026-07-16T17:16:52Z", "completed_at": "2026-07-16T17:18:52Z"}]
        metrics = {"build": {"Instance Type": "c6in.4xlarge", "disk_io_mb": {"min": "44", "max": "590", "avg": "467 MB"}}}
        
        report_str = workflow_monitor.generate_report(run_data, jobs, metrics, "/tmp/logs-dir", "json")
        report = json.loads(report_str)
        
        self.assertEqual(report["run"]["id"], 12345)
        self.assertEqual(report["logs_dir"], "/tmp/logs-dir")
        self.assertEqual(report["jobs"][0]["metrics"]["Instance Type"], "c6in.4xlarge")

    def test_parse_duration_seconds(self):
        start = "2026-07-16T17:16:00Z"
        end = "2026-07-16T17:18:30Z"
        self.assertEqual(workflow_monitor.parse_duration_seconds(start, end), 150)
        self.assertEqual(workflow_monitor.parse_duration_seconds(None, end), 0)

    def test_generate_report_json_includes_seconds_and_slowest_jobs(self):
        run_data = {"id": 12345, "status": "completed", "conclusion": "success", "run_started_at": "2026-07-16T17:00:00Z", "updated_at": "2026-07-16T17:10:00Z"}
        jobs = [
            {"name": "fast_job", "status": "completed", "conclusion": "success", "started_at": "2026-07-16T17:00:00Z", "completed_at": "2026-07-16T17:01:00Z"},
            {"name": "slow_job", "status": "completed", "conclusion": "success", "started_at": "2026-07-16T17:00:00Z", "completed_at": "2026-07-16T17:05:00Z"}
        ]
        report_str = workflow_monitor.generate_report(run_data, jobs, {}, "/tmp/logs", "json")
        report = json.loads(report_str)

        self.assertEqual(report["run"]["runtime_seconds"], 600)
        self.assertEqual(report["run"]["avg_job_duration_seconds"], 180)
        self.assertEqual(report["jobs"][0]["duration_seconds"], 60)
        self.assertEqual(report["jobs"][1]["duration_seconds"], 300)

        # Check slowest_jobs array is pre-sorted descending
        self.assertEqual(len(report["slowest_jobs"]), 2)
        self.assertEqual(report["slowest_jobs"][0]["name"], "slow_job")
        self.assertEqual(report["slowest_jobs"][0]["duration_seconds"], 300)
        self.assertTrue(report["slowest_jobs"][0]["is_outlier"])


    def test_generate_report_markdown_includes_longest_jobs_summary(self):
        run_data = {"id": 12345, "status": "completed", "conclusion": "success", "run_started_at": "2026-07-16T17:00:00Z", "updated_at": "2026-07-16T17:10:00Z"}
        jobs = [
            {"name": "slow_job", "status": "completed", "conclusion": "success", "started_at": "2026-07-16T17:00:00Z", "completed_at": "2026-07-16T17:05:00Z"}
        ]
        report = workflow_monitor.generate_report(run_data, jobs, {}, "/tmp/logs", "markdown")
        self.assertIn("## Longest Jobs (Bottlenecks)", report)
        self.assertIn("slow_job", report)

    def test_generate_report_markdown_links(self):
        run_data = {
            "id": 12345,
            "status": "completed",
            "conclusion": "success",
            "html_url": "https://github.com/owner/repo/actions/runs/12345",
            "run_started_at": "2026-07-16T17:00:00Z",
            "updated_at": "2026-07-16T17:10:00Z"
        }
        jobs = [
            {
                "name": "build",
                "status": "completed",
                "conclusion": "success",
                "html_url": "https://github.com/owner/repo/actions/runs/12345/job/67890",
                "started_at": "2026-07-16T17:00:00Z",
                "completed_at": "2026-07-16T17:05:00Z"
            }
        ]
        report = workflow_monitor.generate_report(run_data, jobs, {}, "/tmp/logs", "markdown")
        self.assertIn("- **ID**: [12345](https://github.com/owner/repo/actions/runs/12345)", report)
        self.assertIn("- **Status**: [completed](https://github.com/owner/repo/actions/runs/12345)", report)
        self.assertIn("| [build](https://github.com/owner/repo/actions/runs/12345/job/67890) |", report)
        self.assertIn("### Job: [build](https://github.com/owner/repo/actions/runs/12345/job/67890)", report)
        self.assertIn("- **Status**: [completed](https://github.com/owner/repo/actions/runs/12345/job/67890)", report)

    @patch('workflow_monitor.parse_args')
    @patch('workflow_monitor.get_trials_dir')
    @patch('workflow_monitor.wait_for_run')
    @patch('workflow_monitor.wait_for_completion')
    @patch('workflow_monitor.download_logs')
    @patch('workflow_monitor.parse_job_logs')
    @patch('workflow_monitor.fetch_jobs')
    @patch('workflow_monitor.generate_report')
    @patch('workflow_monitor.log_stderr')
    @patch('builtins.print')
    def test_main_writes_reports(self, mock_print, mock_log_stderr, mock_gen_report, mock_fetch_jobs, mock_parse_logs, mock_download, mock_wait_complete, mock_wait_run, mock_get_trials_dir, mock_parse_args):
        with tempfile.TemporaryDirectory() as tmp_dir:
            mock_get_trials_dir.return_value = tmp_dir
            args = MagicMock()
            args.run_id = 12345
            args.trial_name = "baseline-12345"
            args.repo = "owner/repo"
            args.token = "fake_token"
            args.poll_interval = 10
            mock_parse_args.return_value = args

            mock_wait_run.return_value = {"id": 12345, "path": ".github/workflows/integration-tests.yml"}
            mock_wait_complete.return_value = {"id": 12345}
            mock_gen_report.return_value = '{"fake": "report"}'
            mock_download.return_value = True
            mock_parse_logs.return_value = {}
            mock_fetch_jobs.return_value = []

            workflow_monitor.main()

            expected_output_dir = os.path.join(tmp_dir, "integration-tests.yml", "baseline-12345")
            expected_logs_dir = os.path.join(expected_output_dir, "logs")
            expected_json = os.path.join(expected_output_dir, "report.json")
            expected_md = os.path.join(expected_output_dir, "report.md")

            self.assertTrue(os.path.isdir(expected_output_dir))
            self.assertTrue(os.path.isdir(expected_logs_dir))
            mock_download.assert_called_once_with("owner/repo", 12345, "fake_token", expected_logs_dir)
            mock_gen_report.assert_any_call({"id": 12345}, [], {}, expected_logs_dir, "json")
            mock_gen_report.assert_any_call({"id": 12345}, [], {}, expected_logs_dir, "markdown")
            self.assertTrue(os.path.exists(expected_json))
            self.assertTrue(os.path.exists(expected_md))
            with open(expected_json, 'r') as f:
                self.assertEqual(f.read(), '{"fake": "report"}')
            with open(expected_md, 'r') as f:
                self.assertEqual(f.read(), '{"fake": "report"}')

    def test_sanitize_dir_name(self):
        self.assertEqual(workflow_monitor.sanitize_dir_name("foo/bar"), "foo_bar")
        self.assertEqual(workflow_monitor.sanitize_dir_name("a\\b"), "a_b")
        self.assertEqual(workflow_monitor.sanitize_dir_name(".."), "_")
        self.assertEqual(workflow_monitor.sanitize_dir_name("  my trial  "), "my trial")

    def test_derive_workflow_name(self):
        self.assertEqual(workflow_monitor.derive_workflow_name({"path": ".github/workflows/integration-tests.yml"}), "integration-tests.yml")
        self.assertEqual(workflow_monitor.derive_workflow_name({"name": "My Workflow"}), "My Workflow")
        self.assertEqual(workflow_monitor.derive_workflow_name({}), "unknown")

    def test_normalize_name(self):
        self.assertEqual(workflow_monitor.normalize_name("build"), "build")
        self.assertEqual(workflow_monitor.normalize_name("build-job"), "buildjob")
        
        name = "Run CCIP v1.6 E2E Tests For Workflow Dispatch / smoke/ccip/ccip_reorg_test.go:GreaterThanFinalityTests"
        expected = "runccipv16e2etestsforworkflowdispatchsmokeccipccipreorgtestgogreaterthanfinalitytests"
        self.assertEqual(workflow_monitor.normalize_name(name), expected)

    def test_matches_job_name(self):
        api_name = "Run CCIP v1.6 E2E Tests For Workflow Dispatch / smoke/ccip/ccip_reorg_test.go:GreaterThanFinalityTests"
        log_job_name = "Run CCIP v1.6 E2E Tests For Workflow Dispatch _ smoke_ccip_ccip_reorg_test.goGreaterThanFinalityTests"
        self.assertTrue(workflow_monitor.matches_job_name(api_name, log_job_name))
        
        # Test suffix or path basename matching
        api_name_path = "Run CCIP v1.6 E2E Tests / smoke/ccip/ccip_reorg_test.go:GreaterThanFinalityTests"
        log_job_name_suffix = "smoke_ccip_ccip_reorg_test.goGreaterThanFinalityTests"
        self.assertTrue(workflow_monitor.matches_job_name(api_name_path, log_job_name_suffix))

if __name__ == '__main__':
    unittest.main()
