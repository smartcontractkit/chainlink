---
name: jira-operations
description: Index of shared JIRA operations for flaky-test skills. Defines the canonical slim-record schema, required env, and how to use these files from another skill.
---

<when_to_use>
These files are **includable references**, not user-facing skills (no `SKILL.md`). Read the relevant file(s) when you need to perform one of the operations listed below. Each file is self-contained: it declares its inputs, steps, and outputs. You do not need to read files you are not using.
</when_to_use>

<requirements>
1. Atlassian MCP available.
2. If MCP is unauthenticated prompt user to authenticate before proceeding.
If Atlassian MCP is not available prompt the user to install it. If authentication fails STOP.
</requirements>

<operation_requirements>
All operations require:
- `cloudId` — the Atlassian cloud ID (from `mcp__atlassian__getAccessibleAtlassianResources`).
- `accountId` — the current user's Atlassian account ID (from `mcp__atlassian__atlassianUserInfo`)
</operation_requirements>

<absolute_constraints>
1. Never return raw JIRA API objects — caller only receives slim records.
</absolute_constraints>

<available_operations>
1. [investigation-comment](./investigation-comment.md) - Comment format for Investigation Updates; parsing prior-attempt comments or adding new comments
2. [fetch-flaky-tickets](./fetch-flaky-tickets.md) | JQL search loop: fetch N eligible flaky-test tickets for a project key
3. [transition-ticket](./transition-ticket.md) - Transition a ticket to a semantic target state
4. [claim-ticket](./claim-ticket.md) - Claim one or more explicitly provided JIRA tickets for work
5. [abandon-ticket](./abandon-ticket.md) - Mid-flight abandonment: unassign → Open → Investigation Update comment
</available_operations>

<canonical_slim_record>
You MUST use [this](./slim-record.md) JSON structure to pass data around in order to avoid calling Atlassian MCP multiple times to read information.
</canonical_slim_record>

<logic>
1. If user provided specific JIRA tickets execute the [claim-ticket](./claim-ticket.md) process.
2. If user asked to work on N eligible tickets execute the [fetch-flaky-tickets](./fetch-flaky-tickets.md) loop.
3. Prepare and return slim records.
4. Once the work on flaky tests has finished, regardless of the result, for each ticket:
    a. Update the ticket with investigation comment
    b. Transition the ticket to correct state — pass `original_assignee` AND `accountId` from the slim record so `transition-ticket.md` step 4 can handle reassignment:
    - abandoned -> `Open`
    - mismatch -> `Open`
    - fixed -> `In Review` (transition-ticket.md step 4 will reassign back to `original_assignee`, or unassign if `original_assignee` is null or matches `accountId`)
</logic>