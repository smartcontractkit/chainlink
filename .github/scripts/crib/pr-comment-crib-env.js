#!/usr/bin/env node

const core = require("@actions/core");
const github = require("@actions/github");
const { Octokit } = require("@octokit/rest");

async function commentExists(octokit, owner, repo, prNumber, uniqueIdentifier) {
  // This will automatically paginate through all comments
  const comments = await octokit.paginate(octokit.rest.issues.listComments, {
    owner,
    repo,
    issue_number: prNumber,
  });

  // Check each comment for the unique identifier
  return comments.some((comment) => comment.body.includes(uniqueIdentifier));
}

async function run() {
  try {
    const token = process.env.GITHUB_TOKEN;
    const octokit = new Octokit({ auth: token });

    const context = github.context;
    const labelsToCheck = ["crib"];
    const { owner, repo } = context.repo;
    const prNumber = context.issue.number;

    if (!prNumber) {
      core.setFailed("Could not get PR number from context");
      return;
    }

    // List labels on the PR
    const { data: labels } = await octokit.rest.issues.listLabelsOnIssue({
      owner,
      repo,
      issue_number: prNumber,
    });

    // Check if any label matches the labelsToCheck
    const labelMatches = labels.some((label) =>
      labelsToCheck.includes(label.name)
    );

    if (!labelMatches) {
      core.info("No 'crib' PR label found. Proceeding.");
      return;
    }

    // Comment header and unique identifier
    const commentHeader = "## CRIB Environment Details";

    // Check if the comment already exists
    if (await commentExists(octokit, owner, repo, prNumber, commentHeader)) {
      core.info("CRIB environment comment already exists. Skipping.");
      return;
    }

    // Construct the comment
    const comment = `${commentHeader} :information_source:

CRIB activated via the 'crib' label. To destroy the environment, remove the 'crib' PR label or close the PR.

Please review the following details:

### Subdomains

_Use these subdomains to access the CRIB environment. They are prefixes to the internal base domain._

- crib-chainlink-${prNumber}-node1.
- crib-chainlink-${prNumber}-node2.
- crib-chainlink-${prNumber}-node3.
- crib-chainlink-${prNumber}-node4.
- crib-chainlink-${prNumber}-node5.
- crib-chainlink-${prNumber}-node6.
- crib-chainlink-${prNumber}-geth-http.
- crib-chainlink-${prNumber}-geth-ws.
- crib-chainlink-${prNumber}-mockserver.
`;

    // Create a comment on the PR
    await octokit.rest.issues.createComment({
      owner,
      repo,
      issue_number: prNumber,
      body: comment,
    });
  } catch (error) {
    core.setFailed(error.message);
  }
}

// Run the script if it's executed directly from the command line
if (require.main === module) {
  run();
}

// Export the run function for testing purposes
module.exports = { run };
