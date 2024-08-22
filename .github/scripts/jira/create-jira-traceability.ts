import * as jira from "jira.js";
import { createJiraClient, extractJiraIssueNumbersFrom } from "./lib";
import * as core from "@actions/core";
import { URL } from "url";

/**
 * Adds traceability to JIRA issues by commenting on each issue with a link to the artifact payload
 * along with a label to connect all issues to the same chainlink product review.
 *
 * @param client The jira client
 * @param issues The list of JIRA issue numbers to add traceability to
 * @param label The label to add to each issue
 * @param artifactUrl The url to the artifact payload that we'll comment on each issue with
 */
export async function addTraceabillityToJiraIssues(
  client: jira.Version3Client,
  issues: string[],
  label: string,
  artifactUrl: string
) {
  for (const issue of issues) {
    await checkAndAddArtifactPayloadComment(client, issue, artifactUrl);

    // CHECK: We don't need to see if the label exists, should no-op
    core.info(`Adding label ${label} to issue ${issue}`);
    await client.issues.editIssue({
      issueIdOrKey: issue,
      update: {
        labels: [{ add: label }],
      },
    });
  }
}

/**
 * Checks if the artifact payload already exists as a comment on the issue, if not, adds it.
 */
async function checkAndAddArtifactPayloadComment(
  client: jira.Version3.Version3Client,
  issue: string,
  artifactUrl: string
) {
  const maxResults = 5000;
  const comments = await client.issueComments.getComments({
    issueIdOrKey: issue,
    maxResults, // this is the default maxResults, see https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-comments/#api-rest-api-3-issue-issueidorkey-comment-get
  });
  if (comments.total ?? 0 > maxResults) {
    throw Error(
      `Too many comments on issue ${issue}, please increase maxResults`
    );
  }
  const commentExists = comments.comments?.some((c) =>
    c?.body?.text?.includes(artifactUrl)
  );

  if (commentExists) {
    core.info(`Artifact payload already exists as comment on issue, skipping`);
  } else {
    core.info(`Adding artifact payload as comment on issue ${issue}`);
    await client.issueComments.addComment({
      issueIdOrKey: issue,
      comment: `:link: [Artifact Payload](${artifactUrl})`,
    });
  }
}

function fetchEnvironmentVariables() {
  const product = process.env.CHAINLINK_PRODUCT;
  if (!product) {
    throw Error("CHAINLINK_PRODUCT environment variable is missing");
  }
  const baseRef = process.env.BASE_REF;
  if (!baseRef) {
    throw Error("BASE_REF environment variable is missing");
  }
  const headRef = process.env.HEAD_REF;
  if (!headRef) {
    throw Error("HEAD_REF environment variable is missing");
  }

  const artifactUrl = process.env.ARTIFACT_URL;
  if (!artifactUrl) {
    throw Error("ARTIFACT_URL environment variable is missing");
  }
  return { product, baseRef, headRef, artifactUrl };
}

function extractChangesetFiles(): string[] {
  const changesetFiles = process.env.CHANGESET_FILES;
  if (!changesetFiles) {
    throw Error("Missing required environment variable CHANGESET_FILES");
  }
  const parsedChangesetFiles = JSON.parse(changesetFiles);
  if (parsedChangesetFiles.length === 0) {
    throw Error("At least one changeset file must exist");
  }

  core.info(`Changeset to extract issues from: ${parsedChangesetFiles.join(", ")}`);
  return parsedChangesetFiles;
}

export function generateJiraIssuesLink(issues: string[]) {
  // https://smartcontract-it.atlassian.net/issues/?jql=issuekey%20in%20%28KS-435%2C%20KS-434%29
  const baseUrl = "https://smartcontract-it.atlassian.net/issues/";
  const jqlQuery = `issuekey in (${issues.join(", ")})`;
  const fullUrl = new URL(baseUrl);
  fullUrl.searchParams.set("jql", jqlQuery);

  const urlStr = fullUrl.toString();
  core.info(`Jira issues link: ${urlStr}`);
  return urlStr
}

export function generateIssueLabel(product: string, baseRef: string, headRef: string) {
  return `review-artifacts-${product}-base:${baseRef}-head:${headRef}`;
}

/**
 * For all affected jira issues listed within the changeset files supplied,
 * we update each jira issue so that they are all labelled and have a comment linking them
 * to the relevant artifact URL.
 */
export async function main() {
  const { product, baseRef, headRef, artifactUrl } =
    fetchEnvironmentVariables();
  const changesetFiles = extractChangesetFiles();
  core.info(`Extracting Jira issue numbers from changeset files: ${changesetFiles.join(", ")}`);
  const jiraIssueNumbers = await extractJiraIssueNumbersFrom(changesetFiles);

  const client = createJiraClient();
  const label = generateIssueLabel(product, baseRef, headRef);
  await addTraceabillityToJiraIssues(
    client,
    jiraIssueNumbers,
    label,
    artifactUrl
  );

  core.setOutput("jira-issues-link", generateJiraIssuesLink(jiraIssueNumbers));
}
