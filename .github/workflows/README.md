# GitHub Workflows

## Required Checks and Path Filters

We want to run certain workflows only when certain file paths change. We can accomplish this with [path based filtering on GitHub actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#onpushpull_requestpull_request_targetpathspaths-ignore). The problem that we run into is that we have certain required checks on GitHub that will not run or pass if we have path based filtering that never executes the workflow.

The [solution that GitHub recommends](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/troubleshooting-required-status-checks#handling-skipped-but-required-checks) is to create a "dummy" workflow with the same workflow name and job names as the required workflow/jobs with the jobs running a command to simply exit zero immediately to indicate success.

### Solution

If your workflow is named `solidity.yml`, create a `solidity-paths-ignore.yml` file with the same workflow name, event triggers (except for the path filters, use `paths-ignore` instead of `paths`), same job names, and then in the steps feel free to echo a command or explicitly `exit 0` to make sure it passes. See the workflow file names with the `-paths-ignore.yml` suffix in this directory for examples.
