import {generateIssueLabel, generateJiraIssuesLink} from './create-jira-traceability';
import {describe, it, expect} from 'vitest'

describe('generateJiraIssuesLink', () => {
  it('should generate a Jira issues link', () => {
    const issues = ['KS-435', 'KS-434'];
    expect(generateJiraIssuesLink(issues)).toMatchInlineSnapshot(`"https://smartcontract-it.atlassian.net/issues/?jql=issuekey+in+%28KS-435%2C+KS-434%29"`)
  });
});

describe('generateIssueLabel', () => {
  it('should generate an issue label', () => {
    const product = 'automation';
    const baseRef = '0de9b3b';
    const headRef = 'e5b3b9d'; 
    expect(generateIssueLabel(product, baseRef, headRef)).toMatchInlineSnapshot(`"review-artifacts-automation-base:0de9b3b-head:e5b3b9d"`)
  });
})
