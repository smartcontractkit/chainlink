import { expect, describe, it } from "vitest";
import { parseIssueNumberFrom, tagsToLabels } from "./lib";

describe("parseIssueNumberFrom", () => {
  it("should return the first JIRA issue number found", () => {
    let r = parseIssueNumberFrom("CORE-123", "CORE-456", "CORE-789");
    expect(r).to.equal("CORE-123");

    r = parseIssueNumberFrom(
      "2f3df5gf",
      "chore/test-RE-78-branch",
      "RE-78 Create new test branches"
    );
    expect(r).to.equal("RE-78");

    // handle lower case
    r = parseIssueNumberFrom("core-123", "CORE-456", "CORE-789");
    expect(r).to.equal("CORE-123");
  });

  it("works with multiline commit bodies", () => {
    const r = parseIssueNumberFrom(
      `This is a multiline commit body

CORE-1011`,
      "CORE-456",
      "CORE-789"
    );
    expect(r).to.equal("CORE-1011");
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
