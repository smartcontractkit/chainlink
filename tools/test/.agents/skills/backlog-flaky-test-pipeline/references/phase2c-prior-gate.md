---
phase: phase2c
model: haiku
---

<phase id="phase2c">

<purpose>
Gate on prior investigation history, collect user decisions for affected tickets, then claim approved tickets. Always entered after phase2a or phase2b. Records gate decisions in `user_decisions` and updates `ticket_records` with approval status.
</purpose>

<prereqs>

Read [../../_shared-jira-flaky-ops/investigation-comment.md](../../_shared-jira-flaky-ops/investigation-comment.md) before parsing any previous-investigation comments. It defines the comment format and parsing rules.
</prereqs>

<step id="1-categorise">
For each slim record in `ticket_records`, classify:

- `previous_attempts` is empty → **auto-approved**, no prompt needed.
- Any attempt with `outcome ≠ "FIXED"` (i.e. `INCONCLUSIVE`, `PARTIAL_FIX`, `MISMATCH`, `SKIP_TOP_LEVEL`, `RETURNED_TO_QUEUE`, `ABANDONED`) → **requires user decision**.

If all tickets are auto-approved: skip step 2, go directly to step 3.
</step>

<step id="2-surface-prior-attempts">
Present each ticket that requires a decision. One prompt per ticket:

<user-prompt>
**KEY-NNN** has {N} prior investigation attempt(s). Most recent ({date}): {outcome}.

**What was tried**: {excluded_approaches joined, or "nothing concrete"}
**Why it failed**: {rejection_reasons joined, or "see comment"}
**Recommended next step**: {recommended_next_step or "none recorded"}

How would you like to proceed?
(a) Continue — use prior findings as context for a new investigation
(b) Skip — return ticket to queue without investigating
</user-prompt>

Apply each response:
- **(a)** → mark approved, record in `user_decisions`.
- **(b)** → mark skipped. Do NOT claim. Do NOT add a JIRA comment — the ticket was never touched, just leave it.

In `--auto` mode: automatically choose (a) for all and log: *"Prior investigation found for KEY-NNN ({outcome}) — proceeding with prior context."*
</step>

<step id="3-claim">
For each approved ticket, follow `_shared-jira-flaky-ops/claim-ticket.md` with `jira_key` and `accountId`.

Announce: "Claimed K issues: [KEY-1, KEY-2, ...]. Proceeding to investigation."

If K = 0 (all skipped): stop. Nothing to investigate.
</step>

<on_complete>
Read [phase3-investigation.md](phase3-investigation.md) and follow its instructions.
</on_complete>

</phase>
