---
phase: phase3a-ii
model_tier: lightweight
---

<substep id="3a-ii" mode="jira">

<purpose>
Gather Trunk and CI-run evidence for a JIRA-mode ticket. Calls `trunk-investigation` first; falls back to the diagnose probe when both Trunk and CI evidence are absent.
</purpose>

<step id="trunk">
Read [phase3-trunk-investigation.md](phase3-trunk-investigation.md) and follow its procedure.

Inputs: `test_case_id`, `test_name`, `title`, `ci_run_url`, `jira_key`, `auto_mode`.

On return, `actionable_facts`, `ci_run_evidence`, `trunk_investigation_status`, and `trunk_analysis_url` are written to the ticket record.
</step>

<step id="diagnose-fallback">
Runs only when `actionable_facts == []` AND `ci_run_evidence == null`.

- `auto_mode = true` → run the diagnose probe without prompting. Inform user: *"No Trunk or CI-run evidence for {jira_key}; auto mode running diagnose locally to gather symptoms."* Read [phase3-diagnose-probe.md](phase3-diagnose-probe.md) with `test_name` (part before first `/`), `package`, `caller_context = jira_key`.
- `auto_mode = false` → ask using `AskUserQuestion`:
  - Question: "{jira_key}: No Trunk or CI-run evidence available. Run `diagnose` locally to gather failure symptoms before code analysis? (~2–5 min, scoped to the package)"
  - Options: ["Yes — run diagnose", "No — proceed with code analysis only"]
  - Yes → Read [phase3-diagnose-probe.md](phase3-diagnose-probe.md) with `test_name` (part before first `/`), `package`, `caller_context = jira_key`.
  - No → leave `actionable_facts = []`, proceed.
</step>

<on_complete>
Read [phase3a-iii-top-level-check.md](phase3a-iii-top-level-check.md) and follow its instructions.
</on_complete>

</substep>
