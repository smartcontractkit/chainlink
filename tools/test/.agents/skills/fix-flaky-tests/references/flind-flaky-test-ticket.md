<find-flaky-test-ticket>
JQL search query that finds a flaky test ticket related to test name [TEST_FUNCTION_NAME] and optionally package that contains it [PACKAGE]. To be used, when user did not provide a JIRA issue.

<requirements>
1. Neither JIRA ticket nor project provided by the user.
2. Current repository — `{owner}/{repo}` extracted from `git remote get-url origin`
</requirements>

<search-for-ticket>
1. Execute JQL query: `Test[Short text]" ~ "[TEST_FUNCTION_NAME]" AND labels = flaky-test AND status = Open`.
2. If you have tests' package add `AND "package[short text]" ~ "[PACKAGE]"`
3. If more than one ticket found return a slim record for each of them without claiming, ask the user which one corresponds
to the test she wants to work on and then claim it.
4. Otherwise claim the ticket by following [claim-ticket](./claim-ticket.md) process.
</search-for-ticket>

</find-flaky-test-ticket>

