import * as core from "@actions/core";
import jira from "jira.js";
import { tagsToLabels, createJiraClient, parseIssueNumberFrom } from "./lib";

function updateJiraIssue(
  client: jira.Version3Client,
  issueNumber: string,
  tags: string[],
  fixVersionName: string,
  dryRun: boolean
) {
  const payload = {
    issueIdOrKey: issueNumber,
    update: {
      labels: tagsToLabels(tags),
      fixVersions: [{ set: [{ name: fixVersionName }] }],
    },
  };

  core.info(
    `Updating JIRA issue ${issueNumber} with fix version ${fixVersionName} and labels [${payload.update.labels.join(
      ", "
    )}]`
  );
  if (dryRun) {
    core.info("Dry run enabled, skipping JIRA issue update");
    return;
  }

  return client.issues.editIssue(payload);
}

async function main() {
  const prTitle = process.env.PR_TITLE;
  const commitMessage = process.env.COMMIT_MESSAGE;
  const branchName = process.env.BRANCH_NAME;

  const chainlinkVersion = process.env.CHAINLINK_VERSION;
  const dryRun = !!process.env.DRY_RUN;
  // tags are not getting used at the current moment so will always default to []
  const tags = process.env.FOUND_TAGS ? process.env.FOUND_TAGS.split(",") : [];

  const client = createJiraClient();

  // Checks for the Jira issue number and exit if it can't find it
  const issueNumber = parseIssueNumberFrom(prTitle, commitMessage, branchName);
  if (!issueNumber) {
    const msg =
      "No JIRA issue number found in: PR title, commit message, or branch name. Please include the issue ID in one of these.";

    core.info(msg);
    core.notice(msg);
    core.setOutput("jiraComment", `> :medal_military: ${msg}`);

    return;
  }

  const fixVersionName = `chainlink-v${chainlinkVersion}`;
  await updateJiraIssue(client, issueNumber, tags, fixVersionName, dryRun);

  core.setOutput("jiraComment", "");
}

async function run() {
  try {
    await main();
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
    core.setFailed(
      "Error: Failed to update JIRA issue with fix version and labels."
    );
  }
}

run();
