# Test Flake Skip Process

When flaky tests are found and are affecting developers ability to merge code we have some steps to take to allow them to be skipped.

1. The owners of the test code should be notified. This will be done automatically at some point in the future.
2. A jira ticket will be created for the flake and assigned to the team that owns the code.
3. The owners of the code can decide whether or not it is too risky to skip the test before the flake is fixed.
4. If the test is to be skipped it should be done using the testing.T.Skip function and provide the jira ticket number for tracking purposes.
