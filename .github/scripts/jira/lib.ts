
import * as core from '@actions/core'
import * as jira from 'jira.js'

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

  core.debug(`Parsing issue number from: ${inputs.join(", ")}`);
  const parsed: string[] = inputs.map(parse).filter((x) => x !== undefined);
  core.debug(`Found issue number: ${parsed[0]}`);

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

export function createJiraClient() {
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
