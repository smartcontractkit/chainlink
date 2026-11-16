from __future__ import annotations

import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent


class BadgeTests(unittest.TestCase):
    def test_badge_uses_source_only_and_writes_state(self):
        with tempfile.TemporaryDirectory() as tmp:
            tmp = Path(tmp)
            report = tmp / "any.json"
            svg = tmp / "badge.svg"
            state = tmp / "report.json"
            report.write_text(
                json.dumps(
                    {
                        "Global": {"Percent": 54.27},
                        "GlobalSourceOnly": {"Percent": 60.92},
                        "Patch": {"Percent": 43.84},
                    }
                ),
                encoding="utf-8",
            )
            subprocess.run(
                [
                    sys.executable,
                    str(HERE / "badge.py"),
                    str(report),
                    str(svg),
                    str(state),
                    "abc123",
                    "coverage-report-abc123-1",
                    "smartcontractkit/chainlink",
                    "https://example.invalid/run/1",
                ],
                check=True,
            )
            svg_text = svg.read_text(encoding="utf-8")
            state_text = json.loads(state.read_text(encoding="utf-8"))
            self.assertIn("coverage: 60.9%", svg_text)
            self.assertIn("60.9%", svg_text)
            self.assertEqual(state_text["global_raw"], 54.27)
            self.assertEqual(state_text["global_source_only"], 60.92)
            self.assertEqual(state_text["patch"], 43.84)
            self.assertEqual(state_text["sha"], "abc123")

    def test_badge_color_changes_with_coverage(self):
        with tempfile.TemporaryDirectory() as tmp:
            tmp = Path(tmp)
            report = tmp / "any.json"
            svg = tmp / "badge.svg"
            state = tmp / "report.json"
            report.write_text(
                json.dumps({"GlobalSourceOnly": {"Percent": 25.0}}),
                encoding="utf-8",
            )
            subprocess.run(
                [
                    sys.executable,
                    str(HERE / "badge.py"),
                    str(report),
                    str(svg),
                    str(state),
                    "abc123",
                    "artifact",
                    "org/repo",
                    "https://example.invalid/run/1",
                ],
                check=True,
            )
            self.assertIn("#c62828", svg.read_text(encoding="utf-8"))


if __name__ == "__main__":
    unittest.main()
