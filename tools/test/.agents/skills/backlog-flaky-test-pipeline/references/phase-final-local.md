---
phase: phase-final-local
model_tier: lightweight
---

<phase id="phase-final-local">

<purpose>
Print the local-mode session summary and end. No commits, no push, no PR, no JIRA writes. Reached only from local mode after phase4 completes.
</purpose>

<summary>
Print the following table, populated from `ticket_records`:

```
Session complete (local mode — no JIRA, no PR).

Fix results:
| Test | Verdict | Local 10x | Notes |
|------|---------|-----------|-------|
| <pkg>.TestFoo  | FIXED        | 10/10 | diff retained, uncommitted |
| <pkg>.TestBar  | PARTIAL_FIX  | 4/10  | reverted |
| <pkg>.TestBaz  | SKIPPED      | —     | classified SUT |
| <pkg>.TestQux  | INCONCLUSIVE | —     | debate did not converge |
```

Column rules:
- **Test**: `{package}.{test_name}` (use `local_id` as fallback if package is null).
- **Verdict**: FIXED | PARTIAL_FIX | SKIPPED | INCONCLUSIVE | MISMATCH.
- **Local 10x**: pass count out of 10 for FIXED/PARTIAL_FIX; `—` for others.
- **Notes**: one short phrase. For FIXED: "diff retained, uncommitted". For PARTIAL_FIX: "reverted". For SKIPPED: the classification reason (e.g. "classified SUT", "classified AMBIGUOUS", "SKIP_TOP_LEVEL"). For INCONCLUSIVE: "debate did not converge". For MISMATCH: "stack trace stale".
</summary>

<footer>
After the table, print:

> FIXED diffs are uncommitted in your working tree. Review with `git diff` and commit manually if you want to keep them.
</footer>

</phase>
