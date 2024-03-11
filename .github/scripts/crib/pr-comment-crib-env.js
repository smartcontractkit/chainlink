#!/usr/bin/env node

import * as core from "@actions/core";
import * as github from "@actions/github";
import { Octokit } from "@octokit/rest";
import { route53RecordsExist } from "./lib/check-route53-records.js";

function generateSubdomains(subdomainPrefix, prNumber) {
  return [
    `${subdomainPrefix}-${prNumber}-node1`,
    `${subdomainPrefix}-${prNumber}-node2`,
    `${subdomainPrefix}-${prNumber}-node3`,
    `${subdomainPrefix}-${prNumber}-node4`,
    `${subdomainPrefix}-${prNumber}-node5`,
    `${subdomainPrefix}-${prNumber}-node6`,
    `${subdomainPrefix}-${prNumber}-geth-http`,
    `${subdomainPrefix}-${prNumber}-geth-ws`,
    `${subdomainPrefix}-${prNumber}-mockserver`,
  ];
}

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
    const route53ZoneId = process.env.ROUTE53_ZONE_ID;

    // Check for the existence of GITHUB_TOKEN and ROUTE53_ZONE_ID
    if (!token || !route53ZoneId) {
      core.setFailed("Error: Missing required environment variables.");
      if (!token) core.setFailed("GITHUB_TOKEN is required.");
      if (!route53ZoneId) core.setFailed("ROUTE53_ZONE_ID is required.");
      return;
    }

    const octokit = new Octokit({ auth: token });
    const context = github.context;
    const subdomainPrefix = process.env.SUBDOMAIN_PREFIX || "crib-chainlink";
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

    // Check if DNS records exist in Route 53 before printing out the subdomains.
    try {
      const subdomains = generateSubdomains(subdomainPrefix, prNumber);
      core.debug("Subdomains:", subdomains);
      const maxRetries = 7; // Maximum number of retries
      const recordsExist = await route53RecordsExist(
        route53ZoneId,
        subdomains,
        maxRetries
      );
      if (recordsExist) {
        core.info("Route 53 DNS records exist:", recordsExist);
      } else {
        core.setFailed(
          "Route 53 DNS records do not exist. Please check the Route 53 hosted zone."
        );
        return;
      }
    } catch (error) {
      core.setFailed(error.message);
      return;
    }

    const subdomainsFormatted = subdomains
      .map((subdomain) => `- ${subdomain}.`)
      .join("\n");

    // Construct the comment
    const comment = `${commentHeader} :information_source:

CRIB activated via the 'crib' label. To destroy the environment, remove the 'crib' PR label or close the PR.

Please review the following details:

### Subdomains

_Use these subdomains to access the CRIB environment. They are prefixes to the internal base domain which work over the VPN._

${subdomainsFormatted}

**NOTE:** If you have trouble resolving these subdomains, please try to reset your VPN DNS and/or local DNS.
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

run();
