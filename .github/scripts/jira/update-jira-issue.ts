import * as core from "@actions/core";
import jira from "jira.js";

/**
 * Given a list of strings, this function will return the first JIRA issue number it finds.
 *
 * @example parseIssueNumberFrom("CORE-123", "CORE-456", "CORE-789") => "CORE-123"
 * @example parseIssueNumberFrom("2f3df5gf", "chore/test-RE-78-branch", "RE-78 Create new test branches") => "RE-78"
 */
export function parseIssueNumberFrom(
  ...inputs: (string | undefined)[]
): string | undefined {
  function parse(str?: string) {
    const jiraIssueRegex = /[A-Z]{2,}-\d+/;

    return str?.toUpperCase().match(jiraIssueRegex)?.[0];
  }

  const parsed: string[] = inputs.map(parse).filter((x) => x !== undefined);

  return parsed[0];
}

/**
 * Converts an array of tags to an array of labels.
 *
 * A label is a string that is formatted as `core-release/{tag}`, with the leading `v` removed from the tag.
 *
 * @example tagsToLabels(["v1.0.0", "v1.1.0"]) => [{ add: "core-release/1.0.0" }, { add: "core-release/1.1.0" }]
 */
export function tagsToLabels(tags: string[]) {
  const labelPrefix = "core-release";

  return tags.map((t) => ({
    add: `${labelPrefix}/${t.substring(1)}`,
  }));
}

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

function createJiraClient() {
  const jiraHost = process.env.JIRA_HOST;
  const jiraUserName = process.env.JIRA_USERNAME;
  const jiraApiToken = process.env.JIRA_API_TOKEN;

  if (!jiraHost || !jiraUserName || !jiraApiToken) {
    core.setFailed(
      "Error: Missing required environment variables: JIRA_HOST and JIRA_USERNAME and JIRA_API_TOKEN."
    );
    process.exit(1);
  }

  return new jira.Version3Client({
    host: jiraHost,
    authentication: {
      basic: {
        email: jiraUserName,
        apiToken: jiraApiToken,
      },
    },
  });
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
