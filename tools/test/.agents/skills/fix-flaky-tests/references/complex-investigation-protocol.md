<complex_investigation_protocol>
A single pass might not be sufficient to understand what causes flakiness in a complex test. In such cases a debate needs to take place.

<roles>
- Proposer
- Challenger
- Arbiter
</roles>

<proposer model="standard">
Reads the test and its helpers analyzing for: timing dependencies, shared global state, ordering assumptions, missing cleanup/teardown, non-deterministic data and other typical flakiness sources. Proposes the most likely root cause and a concrete fix with file and line reference. Returns `proposed_fix_record`.

The investigation must be informed by previously rejected approaches, which are not to be avoided unless there's new conclusive evidence in their favor.

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
Receives the Proposer's full output. Challenges the proposal — alternative causes, edge cases, risk of breaking other tests. Must explicitly take a position on whether the proposed causal mechanism itself is sound. Returns `fix_evaluation_record`.

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

<verdict_record>
```json
{
    "verdict": "PROCEED | INCONCLUSIVE",
    "rationale": "string",
    "next_round": "boolean"
}
```verdict_record
</>
</arbiter>

<absolute_contraints>
1. Up to 3 rounds.
2. Each role is a separate agent call.
3. Never collapse multiple roles into one agent call — self-critique by a single model defeats the purpose.
4. Run per-test investigations in parallel.
</absolute_contraints>

<logic>
1. Gather the evidence from what the user supplied and what `diagnose` tool yielded.
2. Build `discussion_record` and pass it to the Proposer to kick off first round.
3. When Proposer has finished pass `proposed_fix_record  and `discussion_record` to the Challenger.
4. Once Challenger is done pass the `proposed_fix_record`, `fix_evaluation_record` and `discussion_record` to the Arbiter.
5. If Arbiter didn't decide to PROCEED with the proposed fix and round limit wasn't reached:
    a. fill in investigation history
    b. discard proposed fixes and challenges
    c. start new round.
6. Once final verdict is available, build and return `discussion_result_record`.
</logic>

<discussion_record>
```json
{
  "evidence": [
    {
        "source": "user", "ci", "diagnose",
        "content": "string"
    },
  ],
  "investigation_history": {
    "excluded_approaches": ["string"],
    "rejection_reasons": ["string"]
  }
}
```
<field-rules>
1. If Jira comment related to previous fix attempts is present, fill in `investigation_history` when building the record.
</field-rules>
</discussion_record>

<discussion_result_record>
```json
{
  "jira_key": "KEY-NNN | null",
  "test_name: "string",
  "outcome": "PROCEED | INCONCLUSIVE | ABANDONED ",
    "evidence": [
    {
        "source": "user", "ci", "diagnose",
        "content": "string"
    }
  ],
  "fixes": [{
    "fix_file": "string (PROCEED only, null otherwise)",
    "fix_line": "integer (PROCEED only, null otherwise)",
    "fix_description": "string (PROCEED only, null otherwise)",
    }
  ],
  "proposer_root_cause": "string (PROCEED or INCONCLUSIVE only, null otherwise)",
  "excluded_approaches": ["string (PROCEED only, null otherwise)"],
  "fix_description_attempted": "string (INCONCLUSIVE only, null otherwise)",
  "inconclusive_reason": "string (INCONCLUSIVE only, null otherwise)",
  "recommended_next_step": "string | null (INCONCLUSIVE only)"
}
```
</discussion_result_record>
</complex_investigation_protocol>