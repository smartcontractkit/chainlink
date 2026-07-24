import unittest
import json
import os
import tempfile
from unittest.mock import patch, MagicMock
import workflow_compare

class TestWorkflowCompare(unittest.TestCase):

    def test_compare_durations_decrease(self):
        diff, pct_str = workflow_compare.compare_durations("0:03:20", "0:02:40")
        self.assertEqual(diff, "-0:00:40")
        self.assertEqual(pct_str, "-20.0%")

    def test_compare_durations_increase(self):
        diff, pct_str = workflow_compare.compare_durations("0:02:00", "0:02:30")
        self.assertEqual(diff, "+0:00:30")
        self.assertEqual(pct_str, "+25.0%")

    def test_compare_durations_equal(self):
        diff, pct_str = workflow_compare.compare_durations("0:02:00", "0:02:00")
        self.assertEqual(diff, "0:00:00")
        self.assertEqual(pct_str, "0.0%")

    def test_compare_costs_decrease(self):
        diff, pct_str = workflow_compare.compare_costs("$0.1000", "$0.0800")
        self.assertEqual(diff, "-$0.0200")
        self.assertEqual(pct_str, "-20.0%")

    def test_compare_costs_increase(self):
        diff, pct_str = workflow_compare.compare_costs("$0.1000", "$0.1150")
        self.assertEqual(diff, "+$0.0150")
        self.assertEqual(pct_str, "+15.0%")

    def test_compare_costs_equal(self):
        diff, pct_str = workflow_compare.compare_costs("$0.1000", "$0.1000")
        self.assertEqual(diff, "$0.0000")
        self.assertEqual(pct_str, "0.0%")

    def test_compare_metrics(self):
        m1 = {"min": "10.0", "max": "50.0", "avg": "30.0 MB"}
        m2 = {"min": "12.0", "max": "45.0", "avg": "28.5 MB"}
        
        diff = workflow_compare.compare_metrics(m1, m2)
        # We expect it to format average delta: e.g. "avg: 30.0 -> 28.5 MB (-1.5 MB)" or similar
        self.assertIn("avg: 30.0 MB -> 28.5 MB", diff)
        self.assertIn("-1.5 MB", diff)

    def test_generate_comparison(self):
        data1 = {
            "run": {"id": 123, "runtime": "0:10:00", "status": "completed", "conclusion": "success", "total_cost": "$0.1000"},
            "logs_dir": "/tmp/logs1",
            "jobs": [
                {
                    "name": "build",
                    "status": "completed",
                    "conclusion": "success",
                    "runner": {"labels": ["runs-on"], "name": "runner-1"},
                    "duration": "0:03:20",
                    "metrics": {
                        "Instance Type": "c6in.4xlarge",
                        "Cost": "$0.1000",
                        "system.cpu.load_average.1m": {"min": "0.1", "max": "10.0", "avg": "5.0"}
                    }
                }
            ]
        }
        data2 = {
            "run": {"id": 124, "runtime": "0:09:00", "status": "completed", "conclusion": "success", "total_cost": "$0.0800"},
            "logs_dir": "/tmp/logs2",
            "jobs": [
                {
                    "name": "build",
                    "status": "completed",
                    "conclusion": "success",
                    "runner": {"labels": ["runs-on"], "name": "runner-2"},
                    "duration": "0:02:40",
                    "metrics": {
                        "Instance Type": "c7i-flex.8xlarge",
                        "Cost": "$0.0800",
                        "system.cpu.load_average.1m": {"min": "0.2", "max": "12.0", "avg": "6.0"}
                    }
                }
            ]
        }
        
        report = workflow_compare.generate_comparison(data1, data2)
        self.assertIn("# Workflow Trial Comparison", report)
        self.assertIn("c6in.4xlarge", report)
        self.assertIn("c7i-flex.8xlarge", report)
        self.assertIn("- **Runtime Delta**: `-0:01:00` (-10.0%)", report)
        self.assertIn("- **Cost Delta**: `-$0.0200` (-20.0%)", report)
        self.assertIn("Cost: `$0.1000`", report)
        self.assertIn("Cost: `$0.0800`", report)

    def test_generate_comparison_overall_cost_fallback_and_na(self):
        # Fallback when total_cost is not in run dict but in job metrics
        data1 = {
            "run": {"id": 1, "runtime": "0:10:00", "status": "completed", "conclusion": "success"},
            "jobs": [{"name": "j1", "duration": "0:05:00", "metrics": {"Cost": "$0.0500"}}]
        }
        data2 = {
            "run": {"id": 2, "runtime": "0:10:00", "status": "completed", "conclusion": "success"},
            "jobs": [{"name": "j1", "duration": "0:05:00", "metrics": {"Cost": "$0.0700"}}]
        }
        report = workflow_compare.generate_comparison(data1, data2)
        self.assertIn("- **Cost Delta**: `+$0.0200` (+40.0%)", report)

        # N/A when no cost data
        data_no_cost = {
            "run": {"id": 3, "runtime": "0:10:00", "status": "completed", "conclusion": "success"},
            "jobs": [{"name": "j1", "duration": "0:05:00"}]
        }
        report_na = workflow_compare.generate_comparison(data1, data_no_cost)
        self.assertIn("- **Cost Delta**: `N/A` (N/A)", report_na)

    def test_generate_comparison_markdown_links(self):
        data1 = {
            "run": {
                "id": 123,
                "runtime": "0:10:00",
                "status": "completed",
                "conclusion": "success",
                "html_url": "https://github.com/owner/repo/actions/runs/123"
            },
            "jobs": [
                {
                    "name": "build",
                    "status": "completed",
                    "conclusion": "success",
                    "html_url": "https://github.com/owner/repo/actions/runs/123/job/10",
                    "duration": "0:05:00"
                }
            ]
        }
        data2 = {
            "run": {
                "id": 124,
                "runtime": "0:09:00",
                "status": "completed",
                "conclusion": "success",
                "html_url": "https://github.com/owner/repo/actions/runs/124"
            },
            "jobs": [
                {
                    "name": "build",
                    "status": "completed",
                    "conclusion": "success",
                    "html_url": "https://github.com/owner/repo/actions/runs/124/job/20",
                    "duration": "0:04:00"
                }
            ]
        }
        report = workflow_compare.generate_comparison(data1, data2)
        self.assertIn("- **Base Run ID (Trial 1)**: [123](https://github.com/owner/repo/actions/runs/123)", report)
        self.assertIn("Status: [success](https://github.com/owner/repo/actions/runs/123)", report)
        self.assertIn("- **New Run ID (Trial 2)**: [124](https://github.com/owner/repo/actions/runs/124)", report)
        self.assertIn("Status: [success](https://github.com/owner/repo/actions/runs/124)", report)
        self.assertIn("- **Status**: [success](https://github.com/owner/repo/actions/runs/123/job/10) -> [success](https://github.com/owner/repo/actions/runs/124/job/20)", report)

    def test_normalize_name(self):
        self.assertEqual(workflow_compare.normalize_name("build"), "build")
        self.assertEqual(workflow_compare.normalize_name("build-job"), "buildjob")
        
        name = "Run CCIP v1.6 E2E Tests For Workflow Dispatch / smoke/ccip/ccip_reorg_test.go:GreaterThanFinalityTests"
        expected = "runccipv16e2etestsforworkflowdispatchsmokeccipccipreorgtestgogreaterthanfinalitytests"
        self.assertEqual(workflow_compare.normalize_name(name), expected)

    def test_matches_job_name(self):
        api_name = "Run CCIP v1.6 E2E Tests For Workflow Dispatch / smoke/ccip/ccip_reorg_test.go:GreaterThanFinalityTests"
        log_job_name = "Run CCIP v1.6 E2E Tests For Workflow Dispatch _ smoke_ccip_ccip_reorg_test.goGreaterThanFinalityTests"
        self.assertTrue(workflow_compare.matches_job_name(api_name, log_job_name))
        
        # Test suffix or path basename matching
        api_name_path = "Run CCIP v1.6 E2E Tests / smoke/ccip/ccip_reorg_test.go:GreaterThanFinalityTests"
        log_job_name_suffix = "smoke_ccip_ccip_reorg_test.goGreaterThanFinalityTests"
        self.assertTrue(workflow_compare.matches_job_name(api_name_path, log_job_name_suffix))

    def _make_trial_dir(self, trials_dir, workflow, trial, report_data):
        workflow_dir = os.path.join(trials_dir, workflow)
        trial_dir = os.path.join(workflow_dir, trial)
        os.makedirs(trial_dir, exist_ok=True)
        report_path = os.path.join(trial_dir, "report.json")
        with open(report_path, 'w', encoding='utf-8') as f:
            json.dump(report_data, f)
        return report_path

    def test_resolve_trial_by_path(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            trials_dir = tmp_dir
            path = self._make_trial_dir(trials_dir, "wf", "t1", {"run": {"id": 1}})
            workflow, trial, resolved = workflow_compare.resolve_trial(trials_dir, "wf/t1")
            self.assertEqual(workflow, "wf")
            self.assertEqual(trial, "t1")
            self.assertEqual(resolved, path)

    def test_resolve_trial_by_search(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            trials_dir = tmp_dir
            path = self._make_trial_dir(trials_dir, "wf", "t1", {"run": {"id": 1}})
            workflow, trial, resolved = workflow_compare.resolve_trial(trials_dir, "t1")
            self.assertEqual(workflow, "wf")
            self.assertEqual(trial, "t1")
            self.assertEqual(resolved, path)

    def test_resolve_trial_legacy(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            trials_dir = tmp_dir
            workflow_dir = os.path.join(trials_dir, "wf")
            os.makedirs(workflow_dir, exist_ok=True)
            legacy_path = os.path.join(workflow_dir, "t1.json")
            with open(legacy_path, 'w', encoding='utf-8') as f:
                json.dump({"run": {"id": 1}}, f)
            workflow, trial, resolved = workflow_compare.resolve_trial(trials_dir, "t1")
            self.assertEqual(workflow, "wf")
            self.assertEqual(trial, "t1")
            self.assertEqual(resolved, legacy_path)

    def test_resolve_trial_not_found(self):
        with tempfile.TemporaryDirectory() as tmp_dir:
            trials_dir = tmp_dir
            os.makedirs(trials_dir, exist_ok=True)
            with self.assertRaises(FileNotFoundError):
                workflow_compare.resolve_trial(trials_dir, "missing")

    @patch('workflow_compare.parse_args')
    @patch('workflow_compare.get_trials_dir')
    def test_main_creates_comparison(self, mock_get_trials_dir, mock_parse_args):
        with tempfile.TemporaryDirectory() as tmp_dir:
            trials_dir = tmp_dir
            mock_get_trials_dir.return_value = trials_dir

            data1 = {"run": {"id": 100, "runtime": "0:10:00", "status": "completed", "conclusion": "success"}, "jobs": []}
            data2 = {"run": {"id": 101, "runtime": "0:09:00", "status": "completed", "conclusion": "success"}, "jobs": []}
            self._make_trial_dir(trials_dir, "wf", "before", data1)
            self._make_trial_dir(trials_dir, "wf", "after", data2)

            args = MagicMock()
            args.trial_before = "before"
            args.trial_after = "after"
            mock_parse_args.return_value = args

            workflow_compare.main()

            out_file = os.path.join(trials_dir, "wf", "before-after-comparison.md")
            self.assertTrue(os.path.exists(out_file))
            with open(out_file, 'r', encoding='utf-8') as f:
                content = f.read()
            self.assertIn("# Workflow Trial Comparison", content)
            self.assertIn("100", content)
            self.assertIn("101", content)

    @patch('workflow_compare.parse_args')
    @patch('workflow_compare.get_trials_dir')
    def test_main_different_workflows_error(self, mock_get_trials_dir, mock_parse_args):
        with tempfile.TemporaryDirectory() as tmp_dir:
            trials_dir = tmp_dir
            mock_get_trials_dir.return_value = trials_dir
            self._make_trial_dir(trials_dir, "wf1", "t1", {"run": {"id": 1}, "jobs": []})
            self._make_trial_dir(trials_dir, "wf2", "t2", {"run": {"id": 2}, "jobs": []})

            args = MagicMock()
            args.trial_before = "t1"
            args.trial_after = "t2"
            mock_parse_args.return_value = args

            with self.assertRaises(SystemExit) as cm:
                workflow_compare.main()
            self.assertEqual(cm.exception.code, 1)

if __name__ == '__main__':
    unittest.main()
