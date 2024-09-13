import * as core from "@actions/core";
import { createJiraClient, EMPTY_PREFIX, parseIssueNumberFrom, doesIssueExist, PR_PREFIX } from "./lib";
import { appendIssueNumberToChangesetFile, extractChangesetFile } from "./changeset-lib";

async function main() {
  const prTitle = process.env.PR_TITLE;
  const commitMessage = process.env.COMMIT_MESSAGE;
  const branchName = process.env.BRANCH_NAME;
  const dryRun = !!process.env.DRY_RUN;
  const changesetFile = extractChangesetFile();
  const client = createJiraClient();

  // Checks for the Jira issue number and exit if it can't find it
  const issueNumber = parseIssueNumberFrom(EMPTY_PREFIX, prTitle, commitMessage, branchName);
  if (!issueNumber) {
    const msg =
      "No JIRA issue number found in PR title, commit message, or branch name. This pull request must be associated with a JIRA issue.";

    core.setFailed(msg);
    return;
  }

  const exists = await doesIssueExist(client, issueNumber, dryRun);
  if (!exists) {
    core.setFailed(
      `JIRA issue ${issueNumber} not found, this pull request must be associated with a JIRA issue.`
    );
    return;
  }

  core.info(`Appending JIRA issue ${issueNumber} to changeset file`);
  await appendIssueNumberToChangesetFile(PR_PREFIX, changesetFile, issueNumber);
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
