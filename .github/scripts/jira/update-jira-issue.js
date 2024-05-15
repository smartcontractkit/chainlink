#!/usr/bin/env node

import * as core from "@actions/core";
import fetch from "node-fetch";

function parseIssueNumber(prTitle, commitMessage, branchName) {
  const jiraIssueRegex = /[A-Z]{2,}-\d+/;
  if (!!branchName && jiraIssueRegex.test(branchName.toUpperCase())) {
    return branchName.toUpperCase().match(jiraIssueRegex)[0];
  } else if (
    !!commitMessage &&
    jiraIssueRegex.test(commitMessage.toUpperCase())
  ) {
    return commitMessage.toUpperCase().match(jiraIssueRegex)[0];
  } else if (!!prTitle && jiraIssueRegex.test(prTitle.toUpperCase())) {
    return prTitle.toUpperCase().match(jiraIssueRegex)[0];
  } else {
    return null;
  }
}

function getLabels(tags) {
  const labelPrefix = "core-release";
  return tags.map((tag) => {
    return {
      add: `${labelPrefix}/${tag.substring(1)}`,
    };
  });
}

async function updateJiraIssue(
  jiraHost,
  jiraUserName,
  jiraApiToken,
  issueNumber,
  tags,
  fixVersionName
) {
  const token = Buffer.from(`${jiraUserName}:${jiraApiToken}`).toString(
    "base64"
  );
  const bodyData = {
    update: {
      labels: getLabels(tags),
      fixVersions: [{ set: [{ name: fixVersionName }] }],
    },
  };

  fetch(`https://${jiraHost}/rest/api/3/issue/${issueNumber}`, {
    method: "PUT",
    headers: {
      Authorization: `Basic ${token}`,
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(bodyData),
  })
    .then((response) => {
      console.log(`Response: ${JSON.stringify(response)}`);
      return response.text();
    })
    .then((text) => console.log(text))
    .catch((err) => console.error(err));
}

async function run() {
  try {
    const jiraHost = process.env.JIRA_HOST;
    const jiraUserName = process.env.JIRA_USERNAME;
    const jiraApiToken = process.env.JIRA_API_TOKEN;
    const chainlinkVersion = process.env.CHAINLINK_VERSION;
    const prTitle = process.env.PR_TITLE;
    const commitMessage = process.env.COMMIT_MESSAGE;
    const branchName = process.env.BRANCH_NAME;
    // tags are not getting used at the current moment so will always default to []
    const tags = process.env.FOUND_TAGS
      ? process.env.FOUND_TAGS.split(",")
      : [];

    // Check for the existence of JIRA_HOST and JIRA_USERNAME and JIRA_API_TOKEN
    if (!jiraHost || !jiraUserName || !jiraApiToken) {
      core.setFailed(
        "Error: Missing required environment variables: JIRA_HOST and JIRA_USERNAME and JIRA_API_TOKEN."
      );
      return;
    }

    // Checks for the Jira issue number and exit if it can't find it
    const issueNumber = parseIssueNumber(prTitle, commitMessage, branchName);
    if (!issueNumber) {
      core.info(
        "No JIRA issue number found in: PR title, commit message, or branch name. Please include the issue ID in one of these."
      );
      core.notice(
        "No JIRA issue number found in: PR title, commit message, or branch name. Please include the issue ID in one of these."
      );
      return;
    }
    const fixVersionName = `chainlink-v${chainlinkVersion}`;
    await updateJiraIssue(
      jiraHost,
      jiraUserName,
      jiraApiToken,
      issueNumber,
      tags,
      fixVersionName
    );
  } catch (error) {
    core.setFailed(error.message);
  }
}

run();
