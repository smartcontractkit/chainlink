---
name: jira-operations
description: Index of shared JIRA operations for flaky-test skills. Defines the canonical slim-record schema, required env, and how to use these files from another skill.
---

<when_to_use>
These files are **includable references**, not user-facing skills (no `SKILL.md`). Read the relevant file(s) when you need to perform one of the operations listed below. Each file is self-contained: it declares its inputs, steps, and outputs. You do not need to read files you are not using.
</when_to_use>

<requirements>
1. Atlassian MCP available. If Atlassian/Jira MCP unavailable, STOP!. Prompt user to [install it](https://github.com/atlassian/atlassian-mcp-server) before proceeding.
2. If MCP is unauthenticated, STOP! Prompt user to authenticate before proceeding.
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
3. [find-flaky-test-ticket][./find-flaky-test-ticket.md] | JQL search for flaky test ticket by test name
4. [transition-ticket](./transition-ticket.md) - Transition a ticket to a semantic target state
5. [claim-ticket](./claim-ticket.md) - Claim one or more explicitly provided JIRA tickets for work
6. [abandon-ticket](./abandon-ticket.md) - Mid-flight abandonment: unassign → Open → Investigation Update comment
</available_operations>

<canonical_slim_record>
You MUST use [this](./slim-record.md) JSON structure to pass data around in order to avoid calling Atlassian MCP multiple times to read information.
</canonical_slim_record>

<logic>
1. If user provided specific JIRA tickets execute the [claim-ticket](./claim-ticket.md) process.
2. If user asked to work on N eligible tickets execute the [fetch-flaky-tickets](./fetch-flaky-tickets.md) loop.
3. Otherwise try to find a JIRA ticket related to the tests user wants to fix using [find-flaky-test-ticket](./find-flaky-test-ticket) process.
3. Prepare and return slim records.
4. Once the work on flaky tests has finished, regardless of the result, for each ticket:
    a. Update the ticket with investigation comment
    b. Transition the ticket to correct state — always pass `accountId` from `atlassianUserInfo`; pass `original_assignee` from the slim record when abandoning or mismatching:
    - abandoned -> `Open` (unassign via abandon-ticket.md)
    - mismatch -> `Open` (restore `original_assignee` if set)
    - fixed -> `In Review` (transition-ticket.md step 4 assigns to `accountId` — **never unassign** on FIXED)
</logic>

<assignee_rules>
On FIXED, the investigator stays assignee through In Review. Unassigning when `original_assignee == accountId` was incorrect and left tickets with no owner after claim-and-fix.
</assignee_rules>