import { describe, it, expect, vi, beforeAll } from "vitest";
import {
    addChecklistsToIssue, assertChecklistCount,
    cloneIssue, cloneLinkedIssues, copyAllChecklists,
    getChecklistJSONFromIssue,
    getIssueKeys,
    linkIssues,
    transitionIssueWithComment
} from "./enforce-jira-solidity-review";
import * as core from "@actions/core";
import jira from "jira.js";
import axios from "axios";

vi.mock("jira.js");
vi.mock("@actions/core");
vi.mock("axios");

beforeAll(() => {
    process.env.JIRA_HOST = "https://your-jira-host/";
    process.env.JIRA_USERNAME = "your-jira-username";
    process.env.JIRA_API_TOKEN = "your-jira-api-token";
});

describe("getIssueKeys", () => {
    const mockClient = {
        issueSearch: {
            searchForIssuesUsingJql: vi.fn(),
        },
    };

    it("should return issue keys when issues are found", async () => {
        const mockIssues = [{ key: "ISSUE-1" }, { key: "ISSUE-2" }];
        mockClient.issueSearch.searchForIssuesUsingJql.mockResolvedValueOnce({
            issues: mockIssues,
        });

        const result = await getIssueKeys(
            mockClient as unknown as jira.Version3Client,
            "TEST",
            "Bug",
            "summary",
            "Open",
            [],
            10
        );

        expect(result).toEqual(["ISSUE-1", "ISSUE-2"]);
        expect(mockClient.issueSearch.searchForIssuesUsingJql).toHaveBeenCalledWith({
            jql: 'project = TEST AND issuetype = "Bug" AND summary ~ "summary" AND status = "Open"',
            maxResults: 10,
            fields: ["key"],
        });
    });

    it("should return an empty array when no issues are found", async () => {
        mockClient.issueSearch.searchForIssuesUsingJql.mockResolvedValueOnce({
            issues: undefined,
        });

        const result = await getIssueKeys(
            mockClient as unknown as jira.Version3Client,
            "TEST",
            "Bug",
            "summary",
            "Open",
            [],
            10
        );

        expect(result).toEqual([]);
        expect(mockClient.issueSearch.searchForIssuesUsingJql).toHaveBeenCalledWith({
            jql: 'project = TEST AND issuetype = "Bug" AND summary ~ "summary" AND status = "Open"',
            maxResults: 10,
            fields: ["key"],
        });
    });

    it("should handle JQL with issue keys to ignore", async () => {
        const mockIssues = [{ key: "ISSUE-3" }];
        mockClient.issueSearch.searchForIssuesUsingJql.mockResolvedValueOnce({
            issues: mockIssues,
        });

        const result = await getIssueKeys(
            mockClient as unknown as jira.Version3Client,
            "TEST",
            "Bug",
            "summary",
            "Open",
            ["ISSUE-1", "ISSUE-2"],
            10
        );

        expect(result).toEqual(["ISSUE-3"]);
        expect(mockClient.issueSearch.searchForIssuesUsingJql).toHaveBeenCalledWith({
            jql: 'project = TEST AND issuetype = "Bug" AND summary ~ "summary" AND status = "Open" AND issuekey NOT IN (ISSUE-1,ISSUE-2)',
            maxResults: 10,
            fields: ["key"],
        });
    });

    it("should throw an error when the search fails", async () => {
        mockClient.issueSearch.searchForIssuesUsingJql.mockRejectedValueOnce(
            new Error("Search failed")
        );

        await expect(
            getIssueKeys(
                mockClient as unknown as jira.Version3Client,
                "TEST",
                "Bug",
                "summary",
                "Open",
                [],
                10
            )
        ).rejects.toThrow("Search failed");
        expect(mockClient.issueSearch.searchForIssuesUsingJql).toHaveBeenCalledWith({
            jql: 'project = TEST AND issuetype = "Bug" AND summary ~ "summary" AND status = "Open"',
            maxResults: 10,
            fields: ["key"],
        });
    });
});

describe("linkIssues", () => {
    const mockClient = {
        issueLinks: {
            linkIssues: vi.fn(),
        },
    };

    it("should link issues successfully", async () => {
        mockClient.issueLinks.linkIssues.mockResolvedValueOnce(undefined);

        await linkIssues(
            mockClient as unknown as jira.Version3Client,
            "ISSUE-1",
            "ISSUE-2",
            "Blocks"
        );

        expect(mockClient.issueLinks.linkIssues).toHaveBeenCalledWith({
            type: { name: "Blocks" },
            inwardIssue: { key: "ISSUE-1" },
            outwardIssue: { key: "ISSUE-2" },
        });
        expect(core.debug).toHaveBeenCalledWith(
            "Successfully linked issues: ISSUE-1 now 'Blocks' ISSUE-2"
        );
    });

    it("should handle errors when linking issues", async () => {
        const errorMessage = "Linking failed";
        mockClient.issueLinks.linkIssues.mockRejectedValueOnce(new Error(errorMessage));

        await expect(
            linkIssues(
                mockClient as unknown as jira.Version3Client,
                "ISSUE-1",
                "ISSUE-2",
                "Blocks"
            )
        ).rejects.toThrow(errorMessage);

        expect(mockClient.issueLinks.linkIssues).toHaveBeenCalledWith({
            type: { name: "Blocks" },
            inwardIssue: { key: "ISSUE-1" },
            outwardIssue: { key: "ISSUE-2" },
        });
        expect(core.error).toHaveBeenCalledWith(
            `Error linking issues ISSUE-1 and ISSUE-2: Error: ${errorMessage}`
        );
    });
});

describe("cloneIssue", () => {
    const mockClient = {
        issues: {
            getIssue: vi.fn(),
            createIssue: vi.fn(),
        },
    };

    it("should clone an issue successfully", async () => {
        const originalIssue = {
            fields: {
                issuetype: { id: "10001" },
                priority: { id: "2" },
                summary: "Original Issue Summary",
                description: "Original Issue Description",
            },
        };
        const newIssue = { key: "NEW-1" };

        mockClient.issues.getIssue.mockResolvedValueOnce(originalIssue);
        mockClient.issues.createIssue.mockResolvedValueOnce(newIssue);

        const result = await cloneIssue(
            mockClient as unknown as jira.Version3Client,
            "ORIG-1",
            "TEST"
        );

        expect(result).toBe("NEW-1");
        expect(mockClient.issues.getIssue).toHaveBeenCalledWith({
            issueIdOrKey: "ORIG-1",
        });
        expect(mockClient.issues.createIssue).toHaveBeenCalledWith({
            fields: {
                project: { key: "TEST" },
                priority: originalIssue.fields.priority,
                summary: originalIssue.fields.summary,
                description: originalIssue.fields.description,
                issuetype: { id: originalIssue.fields.issuetype.id },
            },
        });
        expect(core.debug).toHaveBeenCalledWith("Cloned issue key: NEW-1");
    });

    it("should throw an error if the original issue is missing issue type id", async () => {
        const originalIssue = {
            fields: {
                issuetype: undefined,
            },
        };

        mockClient.issues.getIssue.mockResolvedValueOnce(originalIssue);

        await expect(
            cloneIssue(
                mockClient as unknown as jira.Version3Client,
                "ORIG-1",
                "TEST"
            )
        ).rejects.toThrow("Issue ORIG-1 is missing issue type id. This should not happen.");

        expect(mockClient.issues.getIssue).toHaveBeenCalledWith({
            issueIdOrKey: "ORIG-1",
        });
        expect(core.error).toHaveBeenCalledWith(
            "Error cloning issue ORIG-1: Error: Issue ORIG-1 is missing issue type id. This should not happen."
        );
    });

    it("should handle errors during cloning", async () => {
        const errorMessage = "Cloning failed";

        mockClient.issues.getIssue.mockRejectedValueOnce(new Error(errorMessage));

        await expect(
            cloneIssue(
                mockClient as unknown as jira.Version3Client,
                "ORIG-1",
                "TEST"
            )
        ).rejects.toThrow(errorMessage);

        expect(mockClient.issues.getIssue).toHaveBeenCalledWith({
            issueIdOrKey: "ORIG-1",
        });
        expect(core.error).toHaveBeenCalledWith(
            `Error cloning issue ORIG-1: Error: ${errorMessage}`
        );
    });
});

describe("transitionIssueWithComment", () => {
    const mockClient = {
        issues: {
            doTransition: vi.fn(),
        },
    };

    it("should transition the issue with the given resolution and comment", async () => {
        mockClient.issues.doTransition.mockResolvedValueOnce(undefined);

        await transitionIssueWithComment(
            mockClient as unknown as jira.Version3Client,
            "ISSUE-1",
            "81",
            "Declined",
            "Closing issue due to an error"
        );

        expect(mockClient.issues.doTransition).toHaveBeenCalledWith({
            issueIdOrKey: "ISSUE-1",
            transition: {
                id: "81",
            },
            fields: {
                resolution: {
                    name: "Declined",
                },
            },
            update: {
                comment: [
                    {
                        add: {
                            body: {
                                type: "doc",
                                version: 1,
                                content: [
                                    {
                                        type: "paragraph",
                                        content: [
                                            {
                                                type: "text",
                                                text: "Closing issue due to an error",
                                            },
                                        ],
                                    },
                                ],
                            },
                        },
                    },
                ],
            },
        });
        expect(core.debug).toHaveBeenCalledWith("Issue ISSUE-1 successfully closed with comment.");
    });

    it("should handle errors during the transition", async () => {
        const errorMessage = "Transition failed";
        mockClient.issues.doTransition.mockRejectedValueOnce(new Error(errorMessage));

        await expect(
            transitionIssueWithComment(
                mockClient as unknown as jira.Version3Client,
                "ISSUE-1",
                "81",
                "Declined",
                "Closing issue due to an error"
            )
        ).rejects.toThrow(errorMessage);

        expect(mockClient.issues.doTransition).toHaveBeenCalledWith({
            issueIdOrKey: "ISSUE-1",
            transition: {
                id: "81",
            },
            fields: {
                resolution: {
                    name: "Declined",
                },
            },
            update: {
                comment: [
                    {
                        add: {
                            body: {
                                type: "doc",
                                version: 1,
                                content: [
                                    {
                                        type: "paragraph",
                                        content: [
                                            {
                                                type: "text",
                                                text: "Closing issue due to an error",
                                            },
                                        ],
                                    },
                                ],
                            },
                        },
                    },
                ],
            },
        });
        expect(core.error).toHaveBeenCalledWith("Failed to update issue ISSUE-1: Error: Transition failed");
    });
});

describe("getChecklistJSONFromIssue", () => {
    const mockChecklist = {
        value: {
            version: "1.0.0",
            checklists: [
                {
                    id: 0,
                    name: "Example to do",
                    items: [
                        {
                            name: "<p>task 1</p>",
                            required: true,
                            completed: true,
                            status: 0,
                            user: "{ user id }",
                            date: "2019-08-13T13:18:24.046Z",
                        },
                        {
                            name: "<p>task 2</p>",
                            required: false,
                            completed: true,
                            status: 0,
                            user: "{ user id }",
                            date: "2019-08-13T13:18:44.988Z",
                        },
                        {
                            name: "<p>task 3</p>",
                            required: false,
                            completed: false,
                            status: 0,
                            user: "{ user id }",
                            date: "2019-08-08T13:30:29.643Z",
                        },
                    ],
                },
            ],
        },
    };

    it("should fetch the checklist JSON from the issue", async () => {
        axios.get.mockResolvedValueOnce({ data: mockChecklist });

        const result = await getChecklistJSONFromIssue("ISSUE-1", "sd-checklists-0");

        expect(result).toEqual(mockChecklist.value);
        expect(axios.get).toHaveBeenCalledWith(
            `https://your-jira-host/rest/api/3/issue/ISSUE-1/properties/sd-checklists-0`,
            {
                auth: {
                    username: process.env.JIRA_USERNAME,
                    password: process.env.JIRA_API_TOKEN,
                },
            }
        );
        expect(core.debug).toHaveBeenCalledWith("Fetching all checklists from issue ISSUE-1");
    });

    it("should handle errors when fetching the checklist JSON", async () => {
        const errorMessage = "Checklist response had unexpected content: {\"some\":\"data\"}";
        axios.get.mockResolvedValueOnce( { data: { some: "data" } });

        await expect(getChecklistJSONFromIssue("ISSUE-1", "sd-checklists-0")).rejects.toThrow(errorMessage);

        expect(axios.get).toHaveBeenCalledWith(
            `https://your-jira-host/rest/api/3/issue/ISSUE-1/properties/sd-checklists-0`,
            {
                auth: {
                    username: process.env.JIRA_USERNAME,
                    password: process.env.JIRA_API_TOKEN,
                },
            }
        );
        expect(core.error).toHaveBeenCalledWith(`Error reading checklists from issue ISSUE-1: Error: ${errorMessage}`);
    });
});

describe("addChecklistsToIssue", () => {
    const issueId = "10001";
    const checklistProperty = "sd-checklists-0";
    const checklistsJson = {
        version: "1.0.0",
        checklists: [
            {
                id: 0,
                name: "Example to do",
                items: [
                    {
                        name: "<p>task 1</p>",
                        required: true,
                        completed: true,
                        status: 0,
                        user: "{ user id }",
                        date: "2019-08-13T13:18:24.046Z",
                    },
                ],
            },
        ],
    };

    it("should add checklists to the issue successfully", async () => {
        axios.put.mockResolvedValueOnce({});

        await addChecklistsToIssue(issueId, checklistProperty, checklistsJson);

        expect(axios.put).toHaveBeenCalledWith(
            `https://your-jira-host/rest/api/3/issue/${issueId}/properties/${checklistProperty}`,
            checklistsJson,
            {
                auth: {
                    username: process.env.JIRA_USERNAME,
                    password: process.env.JIRA_API_TOKEN,
                },
            }
        );
        expect(core.debug).toHaveBeenCalledWith("Added checklists successfully");
    });

    it("should handle errors when adding checklists to the issue", async () => {
        const errorMessage = "Failed to add checklists";
        axios.put.mockRejectedValueOnce(new Error(errorMessage));

        await expect(addChecklistsToIssue(issueId, checklistProperty, checklistsJson)).rejects.toThrow(errorMessage);

        expect(axios.put).toHaveBeenCalledWith(
            `https://your-jira-host/rest/api/3/issue/${issueId}/properties/${checklistProperty}`,
            checklistsJson,
            {
                auth: {
                    username: process.env.JIRA_USERNAME,
                    password: process.env.JIRA_API_TOKEN,
                },
            }
        );
        expect(core.error).toHaveBeenCalledWith(`Failed to add checklists to issue ${issueId}: Error: ${errorMessage}`);
    });
});

describe("assertChecklistCount", () => {
    it("should pass when checklist count is sufficient", () => {
        const checklistJSON = {
            checklists: [
                { id: 0, name: "Checklist 1", items: [] },
                { id: 1, name: "Checklist 2", items: [] },
            ],
        };
        expect(() => assertChecklistCount(checklistJSON, 2)).not.toThrow();
    });

    it("should throw an error when checklist count is insufficient", () => {
        const checklistJSON = {
            checklists: [
                { id: 0, name: "Checklist 1", items: [] },
            ],
        };
        expect(() => assertChecklistCount(checklistJSON, 2)).toThrow("Checklist JSON either did not contain \"checklists\" array or it's lenght was smaller than 2");
    });

    it("should throw an error when checklists key is missing", () => {
        const checklistJSON = {};
        expect(() => assertChecklistCount(checklistJSON, 1)).toThrow("Checklist JSON did not contain any checklists.");
    });
});

const mockChecklistJSON = {
    value: {
        version: "1.0.0",
        checklists: [
            {
                id: 0,
                name: "Example to do",
                items: [
                    {
                        name: "<p>task 1</p>",
                        required: true,
                        completed: true,
                        status: 0,
                        user: "{ user id }",
                        date: "2019-08-13T13:18:24.046Z",
                    },
                ],
            },
        ],
    },
};

describe("copyAllChecklists", () => {
    const sourceIssueId = "10001";
    const targetIssueId = "10002";
    const expectedMinChecklists = 1;


    it("should copy all checklists successfully", async () => {
        axios.get.mockResolvedValueOnce({ data: mockChecklistJSON });
        axios.put.mockResolvedValueOnce({});

        await copyAllChecklists(sourceIssueId, targetIssueId, expectedMinChecklists);

        expect(axios.get).toHaveBeenCalledWith(
            `https://your-jira-host/rest/api/3/issue/${sourceIssueId}/properties/sd-checklists-0`,
            {
                auth: {
                    username: process.env.JIRA_USERNAME,
                    password: process.env.JIRA_API_TOKEN,
                },
            }
        );
        expect(axios.put).toHaveBeenCalledWith(
            `https://your-jira-host/rest/api/3/issue/${targetIssueId}/properties/sd-checklists-0`,
            mockChecklistJSON.value,
            {
                auth: {
                    username: process.env.JIRA_USERNAME,
                    password: process.env.JIRA_API_TOKEN,
                },
            }
        );
        expect(core.debug).toHaveBeenCalledWith(`Copying all checklists from ${sourceIssueId} to ${targetIssueId}`);
    });

    it("should handle errors when copying checklists", async () => {
        axios.get.mockRejectedValueOnce(new Error("Failed to fetch checklists"));

        await expect(copyAllChecklists(sourceIssueId, targetIssueId, expectedMinChecklists)).rejects.toThrow("Failed to fetch checklists");

        expect(core.error).toHaveBeenCalledWith(`Error reading checklists from issue ${sourceIssueId}: Error: Failed to fetch checklists`);
    });
});

const mockJiraClient = {
    issues: {
        getIssue: vi.fn(),
        createIssue: vi.fn(),
    },
    issueLinks: {
        linkIssues: vi.fn(),
    },
};

describe("cloneLinkedIssues", () => {
    const projectKey = "PROJ";
    const sourceIssueKey = "PROJ-1";
    const targetIssueKey = "PROJ-2";
    const issueTypes = ["Task"];
    const expectedLinkedIssues = 1;
    const expectedMinChecklists = 1;

    it("should clone linked issues successfully", async () => {
        const mockLinkedIssue = {
            id: "10001",
            key: "PROJ-3",
            fields: {
                issuetype: { id: "10000", name: "Task" },
                priority: "High",
                summary: "Linked issue summary",
                description: "Linked issue description",
            },
        };
        const mockNewLinkedIssue = {
            key: "PROJ-4",
        };

        mockJiraClient.issues.getIssue.mockResolvedValueOnce({
            fields: {
                issuelinks: [
                    {
                        id: "10001",
                        inwardIssue: mockLinkedIssue,
                    },
                ],
            },
        });
        mockJiraClient.issues.getIssue.mockResolvedValueOnce(mockLinkedIssue);
        mockJiraClient.issues.createIssue.mockResolvedValueOnce(mockNewLinkedIssue);
        mockJiraClient.issueLinks.linkIssues.mockResolvedValueOnce(undefined);

        axios.get.mockResolvedValueOnce({ data: mockChecklistJSON });
        axios.put.mockResolvedValueOnce({});

        await cloneLinkedIssues(mockJiraClient as unknown as jira.Version3Client, projectKey, sourceIssueKey, targetIssueKey, issueTypes, expectedLinkedIssues, expectedMinChecklists);

        expect(mockJiraClient.issues.getIssue).toHaveBeenCalledWith({ issueIdOrKey: sourceIssueKey });
        expect(mockJiraClient.issues.getIssue).toHaveBeenCalledWith({ issueIdOrKey: mockLinkedIssue.key });
        expect(mockJiraClient.issues.createIssue).toHaveBeenCalledWith({
            fields: {
                project: { key: projectKey },
                priority: mockLinkedIssue.fields.priority,
                summary: mockLinkedIssue.fields.summary,
                description: mockLinkedIssue.fields.description,
                issuetype: { id: mockLinkedIssue.fields.issuetype.id },
            },
        });
        expect(mockJiraClient.issueLinks.linkIssues).toHaveBeenCalledWith({
            type: { name: "Blocks" },
            inwardIssue: { key: mockNewLinkedIssue.key },
            outwardIssue: { key: targetIssueKey },
        });
    });

    it("should handle errors when cloning linked issues", async () => {
        mockJiraClient.issues.getIssue.mockRejectedValueOnce(new Error("Failed to fetch issue"));

        await expect(cloneLinkedIssues(mockJiraClient as unknown as jira.Version3Client, projectKey, sourceIssueKey, targetIssueKey, issueTypes, expectedLinkedIssues, expectedMinChecklists)).rejects.toThrow("Failed to fetch issue");

        expect(core.error).toHaveBeenCalledWith(`Error cloning linked issues from ${sourceIssueKey} to ${targetIssueKey}:  Error: Failed to fetch issue`);
    });
});
