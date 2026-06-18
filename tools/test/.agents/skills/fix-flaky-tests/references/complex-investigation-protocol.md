<complex_investigation_protocol>
A single pass might not be sufficient to understand what causes flakiness in a complex test. In such cases a debate needs to take place and additional data points are hightly recommended (CI failure logs, Docker/application logs, etc.).

Before starting the debate ask the user for a CI failure link(s), if not yet provided and analyze them with [github-failure-analyzer](github-failure-analyzer.md) subagent. Accept only Github or other CI providers link, reject Trunk links. Look for stack traces, test logs and workflow artifacts that might contain application logs.

<roles>
<proposer model="standard">
Reads the test and its helpers analyzing for: timing dependencies, shared global state, ordering assumptions, missing cleanup/teardown, non-deterministic data and other typical flakiness sources. Proposes the most likely root cause and a concrete fix with file and line reference. Returns `proposed_fix_record`.

The investigation must be informed by previously rejected approaches, which are not to be avoided unless there's new conclusive evidence in their favor.

**Read budget:** unbounded. Every file/range read MUST be appended to `discussion_record.code_snippets` so later rounds and other roles don't re-read.

<proposed_fix_record>
```json
{
    "root_cause": "string",
    "fixes": [
        {
            "id": "integer",
            "fix_file": "string",
            "fix_line": "integer",
            "fix_description": "string"
        }
    ],
     "excluded_approaches": ["string"]
}
```
</proposed_fix_record>
</proposer>

<challenger model="reasoning">
Receives the Proposer's full output and `discussion_record` (including `code_snippets`). Challenges the proposal — alternative causes, edge cases, risk of breaking other tests. Must explicitly take a position on whether the proposed causal mechanism itself is sound. Returns `fix_evaluation_record`.

**Read budget:** 0 by default. Reasons over `code_snippets`. May perform up to 2 targeted reads to verify a specific claim — each must name `file:start_line-end_line` and the claim it verifies. Any read MUST be appended to `code_snippets`.

<fix_evaluation_record>
```json
[
    {
        "fix_id": "integer",
        "challenges": ["string"]
    }
]
```
</fix_evaluation_record>
</challenger>

<arbiter model="reasoning">
Receives both Proposer and Challenger outputs. Decides whether to stop (enough confidence) or run another round (max 3 total). Issues the final verdict. Returns `verdict_record`.

**Read budget:** 0. Reasons over the records only — no file reads.

<verdict_record>
```json
{
    "decision": "PROCEED | ANOTHER_ROUND | GIVE_UP",
    "rationale": "string"
}
```
</verdict_record>
</arbiter>
</roles>

<absolute_constraints>
1. Up to 3 rounds.
2. Each role is a separate agent call (use the Agent tool).
3. Never collapse multiple roles into one agent call — self-critique by a single model defeats the purpose.
4. Per-test investigations on distinct tickets MAY run in parallel; the 3 roles within a single investigation are sequential.
5. `code_snippets` cache is **per-ticket**. Two tickets MAY share a cache only when one test name is a subtest of the other (e.g. `TestFoo` and `TestFoo/Subcase`) AND the combined cache stays under ~50 entries.
</absolute_constraints>

<logic>
1. Gather evidence from what the user supplied and what `diagnose` yielded. Read prior attempts from `diagnose-attempted-fixes-*.jsonl` and seed `investigation_history` from it. If application or Docker logs are available comb them for errors that might be related to observed test failure.
2. Build `discussion_record` (empty `code_snippets`) and pass it to the Proposer to kick off round 1.
3. When Proposer has finished pass `proposed_fix_record` and the (now-populated) `discussion_record` to the Challenger.
4. Once Challenger is done pass `proposed_fix_record`, `fix_evaluation_record` and `discussion_record` to the Arbiter.
5. Arbiter decision:
    - `PROCEED` → go to step 6.
    - `ANOTHER_ROUND` (only if round < 3): append rejected fix(es) to `discussion_record.investigation_history.excluded_approaches` and Challenger reasoning to `rejection_reasons`. Keep `code_snippets`. Discard `proposed_fix_record` and `fix_evaluation_record`. Start next round.
    - `GIVE_UP` OR round limit reached without PROCEED → outcome is INCONCLUSIVE.
6. Build and return `discussion_result_record`, then follow `<outcome_routing>`.
</logic>

<discussion_record>
```json
{
  "evidence": [
    {
        "source": "user | ci | diagnose",
        "content": "string"
    }
  ],
  "code_snippets": [
    {
        "file": "string",
        "start_line": "integer",
        "end_line": "integer",
        "snippet": "string",
        "why_relevant": "string"
    }
  ],
  "investigation_history": {
    "excluded_approaches": ["string"],
    "rejection_reasons": ["string"]
  }
}
```
<field-rules>
1. If a Jira comment related to previous fix attempts is present, OR a `diagnose-attempted-fixes-*.jsonl` file exists for this test/package, seed `investigation_history` from those sources when building the record.
2. `code_snippets` starts empty; roles append as they read. De-duplicate by `(file, start_line, end_line)`.
</field-rules>
</discussion_record>

<discussion_result_record>
```json
{
  "jira_key": "KEY-NNN | null",
  "test_name": "string",
  "outcome": "PROCEED | INCONCLUSIVE | ABANDONED",
  "evidence": [
    {
        "source": "user | ci | diagnose",
        "content": "string"
    }
  ],
  "fixes": [{
    "fix_file": "string (PROCEED only, null otherwise)",
    "fix_line": "integer (PROCEED only, null otherwise)",
    "fix_description": "string (PROCEED only, null otherwise)"
    }
  ],
  "proposer_root_cause": "string (PROCEED or INCONCLUSIVE only, null otherwise)",
  "excluded_approaches": ["string (PROCEED only, null otherwise)"],
  "inconclusive_reason": "string (INCONCLUSIVE only, null otherwise)",
  "recommended_next_step": "string | null (INCONCLUSIVE only)",
  "abandoned_reason": "string (ABANDONED only, null otherwise)"
}
```
</discussion_result_record>

<outcome_routing>
- `PROCEED` → return to `<loop>` step 7 with `fixes`.
- `INCONCLUSIVE` → write a JSONL entry (`hypothesis` = `proposer_root_cause`, `experiment` = `"complex protocol, N rounds, INCONCLUSIVE"`, `result` = `inconclusive_reason`, `next` = `recommended_next_step`), then:
    - first INCONCLUSIVE on this ticket → re-run `diagnose` with broader signal per `recommended_next_step` (more iterations / `-race` / `-trace`) and re-enter the protocol once. No second extension.
    - else, batch/autonomous mode: run `abandon-ticket.md`.
    - else, interactive mode: ask user [try proposer's guess as a single-pass fix | abandon ticket | redirect with new evidence].
- `ABANDONED` → run `abandon-ticket.md`. Emit this outcome when a role identifies the root cause as 3rd-party-library code (already a STOP per `<absolute_constraints>` in SKILL.md).
</outcome_routing>
</complex_investigation_protocol>