import * as core from "@actions/core";
import jira from "jira.js";
import { createJiraClient, parseIssueNumberFrom } from "./lib";

async function doesIssueExist(
  client: jira.Version3Client,
  issueNumber: string,
  dryRun: boolean
) {
  const payload = {
    issueIdOrKey: issueNumber,
  };

  if (dryRun) {
    core.info("Dry run enabled, skipping JIRA issue enforcement");
    return true;
  }

  try {
    /**
     * The issue is identified by its ID or key, however, if the identifier doesn't match an issue, a case-insensitive search and check for moved issues is performed.
     * If a matching issue is found its details are returned, a 302 or other redirect is not returned. The issue key returned in the response is the key of the issue found.
     */
    const issue = await client.issues.getIssue(payload);
    core.debug(
      `JIRA issue id:${issue.id} key: ${issue.key} found while querying for ${issueNumber}`
    );
    if (issue.key !== issueNumber) {
      core.error(
        `JIRA issue key ${issueNumber} not found, but found issue key ${issue.key} instead. This can happen if the identifier doesn't match an issue, in which case a case-insensitive search and check for moved issues is performed. Make sure the issue key is correct.`
      );
      return false;
    }

    return true;
  } catch (e) {
    core.debug(e as any);
    return false;
  }
}

async function main() {
  const prTitle = process.env.PR_TITLE;
  const commitMessage = process.env.COMMIT_MESSAGE;
  const branchName = process.env.BRANCH_NAME;
  const dryRun = !!process.env.DRY_RUN;
  const client = createJiraClient();

  // Checks for the Jira issue number and exit if it can't find it
  const issueNumber = parseIssueNumberFrom(prTitle, commitMessage, branchName);
  if (!issueNumber) {
    const msg =
      "No JIRA issue number found in PR title, commit message, or branch name. This pull request must be associated with a JIRA issue.";

    core.setFailed(msg);
    return;
  }

  const exists = await doesIssueExist(client, issueNumber, dryRun);
  if (!exists) {
    core.setFailed(`JIRA issue ${issueNumber} not found, this pull request must be associated with a JIRA issue.`);
    return;
  }
}

async function run() {
  try {
    await main();
  } catch (error) {
    if (error instanceof Error) {
      return core.setFailed(error.message);
    }
    core.setFailed(error as any);
  }
}

run();
