import { expect, describe, it, vi } from "vitest";
import {
  generateIssueLabel,
  generateJiraIssuesLink,
  getGitTopLevel,
  parseIssueNumberFrom,
  tagsToLabels,
  EMPTY_PREFIX, SOLIDITY_REVIEW_PREFIX
} from "./lib";
import * as core from "@actions/core";

describe("parseIssueNumberFrom", () => {
  it("should return the first JIRA issue number found", () => {
    let r = parseIssueNumberFrom(EMPTY_PREFIX, "CORE-123", "CORE-456", "CORE-789");
    expect(r).to.equal("CORE-123");

    r = parseIssueNumberFrom(
        EMPTY_PREFIX,
      "chore/test-RE-78-branch",
      "RE-78 Create new test branches"
    );
    expect(r).to.equal("RE-78");

    // handle lower case
    r = parseIssueNumberFrom(EMPTY_PREFIX, "core-123", "CORE-456", "CORE-789");
    expect(r).to.equal("CORE-123");
  });

  it("works with multiline commit bodies", () => {
    let r = parseIssueNumberFrom(
        EMPTY_PREFIX,
      `This is a multiline commit body

CORE-1011`,
      "CORE-456",
      "CORE-789"
    );
    expect(r).to.equal("CORE-1011");

    r = parseIssueNumberFrom(
        SOLIDITY_REVIEW_PREFIX,
        `This is a multiline commit body with prefix

${SOLIDITY_REVIEW_PREFIX}CORE-1011`,
        "CORE-456",
        "CORE-789"
    );
    expect(r).to.equal("CORE-1011");

    r = parseIssueNumberFrom(
        SOLIDITY_REVIEW_PREFIX,
        `This is a multiline commit body with prefix

CORE-111,
${SOLIDITY_REVIEW_PREFIX}CORE-1011`,
        "CORE-456",
        "CORE-789"
    );
    expect(r).to.equal("CORE-1011");
  });

  it("should return undefined if no JIRA issue number is found", () => {
    const result = parseIssueNumberFrom(EMPTY_PREFIX, "No issue number");
    expect(result).to.be.undefined;
  });

  it("works when the label is in the middle of the commit message", () => {
    let r = parseIssueNumberFrom(
        EMPTY_PREFIX,
      "This is a commit message with CORE-123 in the middle",
      "CORE-456",
      "CORE-789"
    );
    expect(r).to.equal("CORE-123");

    r = parseIssueNumberFrom(
        EMPTY_PREFIX,
      "#internal address security vulnerabilities RE-2917 around updating nodes and node operators on capabilities registry"
    );
    expect(r).to.equal("RE-2917");
  });

  it("work also with a prefix", () => {
    let r = parseIssueNumberFrom("PR: ", "PR: RE-78 Create new test branches");
    expect(r).to.equal("RE-78");

    r = parseIssueNumberFrom("PR: ", "RE-99, PR: RE-78 Create new test branches");
    expect(r).to.equal("RE-78");
  });
});

describe("tagsToLabels", () => {
  it("should convert an array of tags to an array of labels", () => {
    const tags = ["v1.0.0", "v1.1.0"];
    const result = tagsToLabels(tags);
    expect(result).to.deep.equal([
      { add: "core-release/1.0.0" },
      { add: "core-release/1.1.0" },
    ]);
  });
});

const mockExecPromise = vi.fn();
vi.mock("util", () => ({
  promisify: () => mockExecPromise,
}));

describe("getGitTopLevel", () => {
  it("should log the top-level directory when git command succeeds", async () => {
    mockExecPromise.mockResolvedValueOnce({
      stdout: "/path/to/top-level-dir",
      stderr: "",
    });

    const mockConsoleLog = vi.spyOn(core, "info");
    await getGitTopLevel();

    expect(mockExecPromise).toHaveBeenCalledWith(
      "git rev-parse --show-toplevel"
    );
    expect(mockConsoleLog).toHaveBeenCalledWith(
      "Top-level directory: /path/to/top-level-dir"
    );
  });

  it("should log an error message when git command fails", async () => {
    mockExecPromise.mockRejectedValueOnce({
      message: "Command failed",
    });

    const mockConsoleError = vi.spyOn(core, "error");
    await getGitTopLevel().catch(() => {});

    expect(mockExecPromise).toHaveBeenCalledWith(
      "git rev-parse --show-toplevel"
    );
    expect(mockConsoleError).toHaveBeenCalledWith(
      "Error executing command: Command failed"
    );
  });

  it("should log an error message when git command output contains an error", async () => {
    mockExecPromise.mockResolvedValueOnce({
      stdout: "",
      stderr: "Error: Command failed",
    });

    const mockConsoleError = vi.spyOn(core, "error");
    await getGitTopLevel().catch(() => {});

    expect(mockExecPromise).toHaveBeenCalledWith(
      "git rev-parse --show-toplevel"
    );
    expect(mockConsoleError).toHaveBeenCalledWith(
      "Error in command output: Error: Command failed"
    );
  });
});

describe("generateJiraIssuesLink", () => {
  it("should generate a Jira issues link", () => {
    expect(
      generateJiraIssuesLink(
        "https://smartcontract-it.atlassian.net/issues/",
        "review-artifacts-automation-base:0de9b3b-head:e5b3b9d"
      )
    ).toMatchInlineSnapshot(
      `"https://smartcontract-it.atlassian.net/issues/?jql=labels+%3D+%22review-artifacts-automation-base%3A0de9b3b-head%3Ae5b3b9d%22"`
    );
  });
});

describe("generateIssueLabel", () => {
  it("should generate an issue label", () => {
    const product = "automation";
    const baseRef = "0de9b3b";
    const headRef = "e5b3b9d";
    expect(generateIssueLabel(product, baseRef, headRef)).toMatchInlineSnapshot(
      `"review-artifacts-automation-base:0de9b3b-head:e5b3b9d"`
    );
  });
});
