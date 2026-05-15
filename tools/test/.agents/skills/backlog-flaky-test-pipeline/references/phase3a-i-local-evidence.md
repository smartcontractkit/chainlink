---
phase: phase3a-i
model_tier: lightweight
---

<substep id="3a-i" mode="local">

<purpose>
Gather observational evidence for a local-mode ticket. Uses the user-supplied log if provided; otherwise runs the diagnose probe to collect failure symptoms before code analysis.
</purpose>

<steps>
- If `slim_record.provided_log_text` is non-null → set `actionable_facts = [slim_record.provided_log_text]`, `trunk_investigation_status = "user_provided"`.
- If `provided_log_text` is null → Read [phase3-diagnose-probe.md](phase3-diagnose-probe.md) and follow its procedure with `test_name` (part before first `/`), `package`, `caller_context = slim_record.test_name`.
</steps>

<on_complete>
Read [phase3a-iii-top-level-check.md](phase3a-iii-top-level-check.md) and follow its instructions.
</on_complete>

</substep>
