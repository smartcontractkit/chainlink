import { readFile } from "fs/promises";
import * as core from "@actions/core";
import * as jira from "jira.js";
import { exec } from "child_process";
import { promisify } from "util";
import { join } from "path";

export async function getGitTopLevel(): Promise<string> {
  const execPromise = promisify(exec);
  try {
    const { stdout, stderr } = await execPromise(
      "git rev-parse --show-toplevel"
    );

    if (stderr) {
      const msg = `Error in command output: ${stderr}`;
      core.error(msg);
      throw Error(msg);
    }

    const topLevelDir = stdout.trim();
    core.info(`Top-level directory: ${topLevelDir}`);
    return topLevelDir;
  } catch (error) {
    const msg = `Error executing command: ${(error as any).message}`;
    core.error(msg);
    throw Error(msg);
  }
}

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

export async function extractJiraIssueNumbersFrom(filePaths: string[]) {
  const issueNumbers: string[] = [];
  const gitTopLevel = await getGitTopLevel();
  for (const path of filePaths) {
    const fullPath = join(gitTopLevel, path);
    const content = await readFile(fullPath, "utf-8");
    const issueNumber = parseIssueNumberFrom(content);
    if (issueNumber) {
      issueNumbers.push(issueNumber);
    }
  }

  return issueNumbers;
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
