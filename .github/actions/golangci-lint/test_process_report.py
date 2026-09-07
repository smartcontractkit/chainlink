#!/usr/bin/env python3
import json
import tempfile
import unittest
from pathlib import Path

import process_report


class TestProcessReport(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.TemporaryDirectory()
        self.report_file = Path(self.temp_dir.name) / "report.json"
        self.output_file = Path(self.temp_dir.name) / "output.txt"
        self.summary_file = Path(self.temp_dir.name) / "summary.md"

    def tearDown(self):
        self.temp_dir.cleanup()

    def test_process_report_with_issues(self):
        data = {
            "Issues": [
                {
                    "FromLinter": "paralleltest",
                    "Text": "Function TestFoo missing parallel\n",
                    "Severity": "high",
                    "Pos": {
                        "Filename": "core/foo_test.go",
                        "Line": 10,
                        "Column": 2,
                    },
                },
                {
                    "FromLinter": "revive",
                    "Text": "exported method Bar should have comment | doc",
                    "Severity": "low",
                    "Pos": {
                        "Filename": "core/foo_test.go",
                        "Line": 25,
                        "Column": 1,
                    },
                },
                {
                    "FromLinter": "revive",
                    "Text": "unused param x",
                    "Severity": "medium",
                    "Pos": {
                        "Filename": "core/bar.go",
                        "Line": 5,
                        "Column": 8,
                    },
                },
            ]
        }
        with open(self.report_file, "w") as f:
            json.dump(data, f)

        summary, outputs = process_report.process(
            report_path=str(self.report_file),
            module_name="core",
            summary_path=str(self.summary_file),
            output_path=str(self.output_file),
        )

        # Check outputs
        self.assertEqual(outputs["file-count"], 2)
        self.assertEqual(outputs["issue-count"], 3)
        self.assertEqual(outputs["critical-issue-count"], 2)  # high (1) + medium (1)
        self.assertIn("2 revive", outputs["per-source-count"])
        self.assertIn("1 paralleltest", outputs["per-source-count"])
        self.assertIn("1 high", outputs["per-severity-count"])
        self.assertIn("1 medium", outputs["per-severity-count"])
        self.assertIn("1 low", outputs["per-severity-count"])

        # Check summary file content
        with open(self.summary_file, "r") as f:
            summary_content = f.read()
        self.assertIn("### 🔍 GolangCI Lint Results: `core`", summary_content)
        self.assertIn("| **Total Issues** | 3 |", summary_content)
        self.assertIn("| **Files Affected** | 2 |", summary_content)
        self.assertIn("| **Critical Issues** | 2 |", summary_content)
        self.assertIn("`core/foo_test.go`", summary_content)
        self.assertIn("exported method Bar should have comment &#124; doc", summary_content)

    def test_process_report_empty_issues(self):
        data = {"Issues": []}
        with open(self.report_file, "w") as f:
            json.dump(data, f)

        summary, outputs = process_report.process(
            report_path=str(self.report_file),
            module_name="plugins",
            summary_path=str(self.summary_file),
            output_path=str(self.output_file),
        )

        self.assertEqual(outputs["file-count"], 0)
        self.assertEqual(outputs["issue-count"], 0)
        self.assertEqual(outputs["critical-issue-count"], 0)

    def test_process_report_missing_file_treated_as_empty(self):
        missing = Path(self.temp_dir.name) / "does-not-exist.json"

        summary, outputs = process_report.process(
            report_path=str(missing),
            module_name="core/scripts",
            summary_path=str(self.summary_file),
            output_path=str(self.output_file),
        )

        self.assertEqual(outputs["file-count"], 0)
        self.assertEqual(outputs["issue-count"], 0)
        self.assertEqual(outputs["critical-issue-count"], 0)
        self.assertEqual(outputs["suffix"], "core-scripts")

        with open(self.output_file, "r") as f:
            self.assertIn("suffix=core-scripts", f.read())

    def test_process_report_custom_critical_severities(self):
        data = {
            "Issues": [
                {
                    "FromLinter": "paralleltest",
                    "Text": "Function TestFoo missing parallel\n",
                    "Severity": "high",
                    "Pos": {"Filename": "core/foo_test.go", "Line": 10, "Column": 2},
                },
                {
                    "FromLinter": "revive",
                    "Text": "exported method Bar should have comment | doc",
                    "Severity": "low",
                    "Pos": {"Filename": "core/foo_test.go", "Line": 25, "Column": 1},
                },
                {
                    "FromLinter": "revive",
                    "Text": "unused param x",
                    "Severity": "medium",
                    "Pos": {"Filename": "core/bar.go", "Line": 5, "Column": 8},
                },
            ]
        }
        with open(self.report_file, "w") as f:
            json.dump(data, f)

        # When only 'high' is critical
        _, outputs = process_report.process(
            report_path=str(self.report_file),
            module_name="core",
            summary_path=str(self.summary_file),
            output_path=str(self.output_file),
            critical_severities=("high",),
        )
        self.assertEqual(outputs["critical-issue-count"], 1)

        # When none are critical
        _, outputs = process_report.process(
            report_path=str(self.report_file),
            module_name="core",
            summary_path=str(self.summary_file),
            output_path=str(self.output_file),
            critical_severities=(),
        )
        self.assertEqual(outputs["critical-issue-count"], 0)
        self.assertEqual(outputs["suffix"], "core")

    def test_get_artifact_suffix(self):
        self.assertEqual(process_report.get_artifact_suffix(""), "root")
        self.assertEqual(process_report.get_artifact_suffix("."), "root")
        self.assertEqual(process_report.get_artifact_suffix("./"), "root")
        self.assertEqual(process_report.get_artifact_suffix("core/scripts"), "core-scripts")
        self.assertEqual(process_report.get_artifact_suffix("core/scripts/"), "core-scripts")
        self.assertEqual(process_report.get_artifact_suffix("plugins"), "plugins")


if __name__ == "__main__":
    unittest.main()

