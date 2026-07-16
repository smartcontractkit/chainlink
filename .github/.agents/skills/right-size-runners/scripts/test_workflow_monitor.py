import unittest
from unittest.mock import patch, MagicMock
import io
import os
import tempfile
import zipfile
import json
import urllib.error
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

    @patch('workflow_monitor.parse_args')
    @patch('workflow_monitor.wait_for_run')
    @patch('workflow_monitor.wait_for_completion')
    @patch('workflow_monitor.download_logs')
    @patch('workflow_monitor.parse_job_logs')
    @patch('workflow_monitor.fetch_jobs')
    @patch('workflow_monitor.generate_report')
    @patch('workflow_monitor.log_stderr')
    @patch('builtins.print')
    def test_main_with_outfile(self, mock_print, mock_log_stderr, mock_gen_report, mock_fetch_jobs, mock_parse_logs, mock_download, mock_wait_complete, mock_wait_run, mock_parse_args):
        args = MagicMock()
        args.run_id = 12345
        args.repo = "owner/repo"
        args.token = "fake_token"
        args.poll_interval = 10
        args.format = "json"
        
        with tempfile.TemporaryDirectory() as tmp_dir:
            out_file = os.path.join(tmp_dir, "report.json")
            args.out_file = out_file
            mock_parse_args.return_value = args
            mock_gen_report.return_value = '{"fake": "report"}'
            
            workflow_monitor.main()
            
            mock_print.assert_not_called()
            self.assertTrue(os.path.exists(out_file))
            with open(out_file, 'r') as f:
                self.assertEqual(f.read(), '{"fake": "report"}')
                
            mock_log_stderr.assert_any_call(f"Report successfully saved to: {out_file}")
            mock_log_stderr.assert_any_call(f"Try exploring with: jq . {out_file}")

    @patch('workflow_monitor.parse_args')
    @patch('workflow_monitor.wait_for_run')
    @patch('workflow_monitor.wait_for_completion')
    @patch('workflow_monitor.download_logs')
    @patch('workflow_monitor.parse_job_logs')
    @patch('workflow_monitor.fetch_jobs')
    @patch('workflow_monitor.generate_report')
    @patch('workflow_monitor.log_stderr')
    @patch('builtins.print')
    def test_main_without_outfile(self, mock_print, mock_log_stderr, mock_gen_report, mock_fetch_jobs, mock_parse_logs, mock_download, mock_wait_complete, mock_wait_run, mock_parse_args):
        args = MagicMock()
        args.run_id = 12345
        args.repo = "owner/repo"
        args.token = "fake_token"
        args.poll_interval = 10
        args.format = "json"
        args.out_file = None
        mock_parse_args.return_value = args
        
        mock_gen_report.return_value = '{"fake": "report"}'
        
        workflow_monitor.main()
        
        mock_print.assert_called_once_with('{"fake": "report"}')

if __name__ == '__main__':
    unittest.main()
