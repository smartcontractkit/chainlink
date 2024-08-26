import * as jira from "jira.js";
import {
  createJiraClient,
  extractJiraIssueNumbersFrom,
  generateIssueLabel,
  generateJiraIssuesLink,
  getJiraEnvVars,
} from "./lib";
import * as core from "@actions/core";

/**
 * Extracts the list of changeset files. Intended to be used with https://github.com/dorny/paths-filter with
 * the 'csv' output format.
 *
 * @returns An array of strings representing the changeset files.
 * @throws {Error} If the required environment variable CHANGESET_FILES is missing.
 * @throws {Error} If no changeset file exists.
 */
function extractChangesetFiles(): string[] {
  const changesetFiles = process.env.CHANGESET_FILES;
  if (!changesetFiles) {
    throw Error("Missing required environment variable CHANGESET_FILES");
  }
  const parsedChangesetFiles = changesetFiles.split(",");
  if (parsedChangesetFiles.length === 0) {
    throw Error("At least one changeset file must exist");
  }

  core.info(
    `Changeset to extract issues from: ${parsedChangesetFiles.join(", ")}`
  );
  return parsedChangesetFiles;
}

/**
 * Adds traceability to JIRA issues by commenting on each issue with a link to the artifact payload
 * along with a label to connect all issues to the same chainlink product review.
 *
 * @param client The jira client
 * @param issues The list of JIRA issue numbers to add traceability to
 * @param label The label to add to each issue
 * @param artifactUrl The url to the artifact payload that we'll comment on each issue with
 */
async function addTraceabillityToJiraIssues(
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
  const getCommentsResponse = await client.issueComments.getComments({
    issueIdOrKey: issue,
    maxResults, // this is the default maxResults, see https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-comments/#api-rest-api-3-issue-issueidorkey-comment-get
  });
  core.debug(JSON.stringify(getCommentsResponse.comments));
  if ((getCommentsResponse.total ?? 0) > maxResults) {
    throw Error(
      `Too many (${getCommentsResponse.total}) comments on issue ${issue}, please increase maxResults (${maxResults})`
    );
  }

  // Search path is getCommentsResponse.comments[].body.content[].content[].marks[].attrs.href
  //
  // Example:
  // [ // getCommentsResponse.comments
  //   {
  //     body: {
  //       type: "doc",
  //       version: 1,
  //       content: [
  //         {
  //           type: "paragraph",
  //           content: [
  //             {
  //               type: "text",
  //               text: "Artifact URL",
  //               marks: [
  //                 {
  //                   type: "link",
  //                   attrs: {
  //                     href: "https://github.com/smartcontractkit/chainlink/actions/runs/10517121836/artifacts/1844867108",
  //                   },
  //                 },
  //               ],
  //             },
  //           ],
  //         },
  //       ],
  //     },
  //   },
  // ];
  const commentExists = getCommentsResponse.comments?.some((c) =>
    c?.body?.content?.some((innerContent) =>
      innerContent?.content?.some((c) =>
        c.marks?.some((m) => m.attrs?.href === artifactUrl)
      )
    )
  );

  if (commentExists) {
    core.info(`Artifact payload already exists as comment on issue, skipping`);
  } else {
    core.info(`Adding artifact payload as comment on issue ${issue}`);
    await client.issueComments.addComment({
      issueIdOrKey: issue,
      comment: {
        type: "doc",
        version: 1,
        content: [
          {
            type: "paragraph",
            content: [
              {
                type: "text",
                text: "Artifact Download URL",
                marks: [
                  {
                    type: "link",
                    attrs: {
                      href: artifactUrl,
                    },
                  },
                ],
              },
            ],
          },
        ],
      },
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

/**
 * For all affected jira issues listed within the changeset files supplied,
 * we update each jira issue so that they are all labelled and have a comment linking them
 * to the relevant artifact URL.
 */
async function main() {
  const { product, baseRef, headRef, artifactUrl } =
    fetchEnvironmentVariables();
  const changesetFiles = extractChangesetFiles();
  core.info(
    `Extracting Jira issue numbers from changeset files: ${changesetFiles.join(
      ", "
    )}`
  );
  const jiraIssueNumbers = await extractJiraIssueNumbersFrom(changesetFiles);

  const client = createJiraClient();
  const label = generateIssueLabel(product, baseRef, headRef);
  await addTraceabillityToJiraIssues(
    client,
    jiraIssueNumbers,
    label,
    artifactUrl
  );

  const { jiraHost } = getJiraEnvVars()
  core.summary.addLink("Jira Issues", generateJiraIssuesLink(`${jiraHost}/issues/`, label));
  core.summary.write();
}
main();
