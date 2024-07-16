import { expect, describe, it } from "vitest";
import { parseIssueNumberFrom, tagsToLabels } from "./update-jira-issue";

describe("parseIssueNumberFrom", () => {
  it("should return the first JIRA issue number found", () => {
    const result1 = parseIssueNumberFrom("CORE-123", "CORE-456", "CORE-789");
    expect(result1).to.equal("CORE-123");

    const result2 = parseIssueNumberFrom(
      "2f3df5gf",
      "chore/test-RE-78-branch",
      "RE-78 Create new test branches"
    );
    expect(result2).to.equal("RE-78");
  });

  it("should return undefined if no JIRA issue number is found", () => {
    const result = parseIssueNumberFrom("No issue number");
    expect(result).to.be.undefined;
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
