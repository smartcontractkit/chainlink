import unittest
import json
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
            "run": {"id": 123, "runtime": "0:10:00", "status": "completed", "conclusion": "success"},
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
            "run": {"id": 124, "runtime": "0:09:00", "status": "completed", "conclusion": "success"},
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
        self.assertIn("-0:00:40", report) # duration delta
        self.assertIn("-$0.0200", report) # cost delta

if __name__ == '__main__':
    unittest.main()
