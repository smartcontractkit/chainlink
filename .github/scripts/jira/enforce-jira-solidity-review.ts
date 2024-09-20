import * as core from "@actions/core";
import jira from "jira.js";
import axios from "axios";
import { join } from "path";
import { createJiraClient, extractJiraIssueNumbersFrom, getJiraEnvVars, doesIssueExist, PR_PREFIX, SOLIDITY_REVIEW_PREFIX } from "./lib";
import { appendIssueNumberToChangesetFile, extractChangesetFile } from "./changeset-lib";
import fs from "fs";

async function main() {
    core.info('Started linking PR to a Solidity Review issue')
    const solidityReviewTemplateKey = readSolidityReviewTemplateKey()
    const changesetFile = extractChangesetFile();

    const jiraPRIssueKeys = await extractJiraIssueNumbersFrom(PR_PREFIX, [changesetFile])
    if (jiraPRIssueKeys.length !== 1) {
        core.setFailed(
            `Solidity Review enforcement only works with 1 JIRA issue per PR, but found ${jiraPRIssueKeys.length} issues in changeset file ${changesetFile}`
          );

          return
    }

    const jiraPRIssueKey = jiraPRIssueKeys[0]
    const client = createJiraClient();

    const jiraSolidityIssues = await extractJiraIssueNumbersFrom(SOLIDITY_REVIEW_PREFIX, [changesetFile])
    if (jiraSolidityIssues.length > 0) {
        for (const jiraIssue of jiraSolidityIssues) {
          const exists = await doesIssueExist(client, jiraIssue, false);
          if (!exists) {
            core.setFailed(
              `JIRA issue ${jiraIssue} not found, Solidity Review issue must exist in JIRA.`
            );

            return;
          }
        }

        core.info(`Found linked Solidity Review issue(s): ${join(...jiraSolidityIssues)}. Nothing more needs to be done.`)

        return
    }

    const jiraProject = extracProjectFromIssueKey(jiraPRIssueKey)
    if (!jiraProject) {
        core.setFailed(`Could not extract project key from issue: ${jiraPRIssueKey}`)

        return
    }

    core.info(`Checking if project ${jiraProject} has any open Solidity Review issues`)
    let solidityReviewIssueKey: string
    const openSolidityReviewIssues = await getOpenSolidityReviewIssuesForProject(client, jiraProject, [solidityReviewTemplateKey])
    if (openSolidityReviewIssues.length === 0) {
        solidityReviewIssueKey = await createSolidityReviewIssue(client, jiraProject, solidityReviewTemplateKey)
    } else if (openSolidityReviewIssues.length === 1) {
        solidityReviewIssueKey = openSolidityReviewIssues[0]
    } else {
        core.setFailed(`Found following open Solidity Review issues for project ${jiraProject}: ${join(...openSolidityReviewIssues)}.
Since we are unable to automatically determine, which one to use, please manualy add it to changeset file: ${changesetFile}. Use this exact format:
${SOLIDITY_REVIEW_PREFIX}<issue-key>
Exmaple with issue key 'PROJ-1234':
${SOLIDITY_REVIEW_PREFIX}PROJ-1234`)

      return
    }

    core.info(`Will use Solidity Review issue: ${solidityReviewIssueKey}`)
    await linkIssues(client, solidityReviewIssueKey, jiraPRIssueKey, 'Blocks')

    core.info(`Appending JIRA Solidity Review issue ${solidityReviewIssueKey} to changeset file`);
    await appendIssueNumberToChangesetFile(SOLIDITY_REVIEW_PREFIX, changesetFile, solidityReviewIssueKey);
    exportIssueKeysToGithubEnv(jiraPRIssueKey, solidityReviewIssueKey)

    core.info('Finished linking PR to a Solidity Review issue')
  }

  /**
   * Retrieves issue keys from Jira based on the specified criteria.
   *
   * This function searches for issues within a given project that match the specified issue type, summary, and status.
   * It can also exclude certain issue keys from the search results and limit the number of results returned.
   *
   * @param client - The Jira client instance used to perform the search.
   * @param projectKey - The key of the project to search within.
   * @param issueType - The type of issues to search for (e.g., 'Bug', 'Task').
   * @param summary - A substring to match within the issue summaries.
   * @param status - The status of the issues to search for (e.g., 'Open', 'Closed').
   * @param issueKeysToIgnore - An array of issue keys to exclude from the search results.
   * @param maxResults - The maximum number of issue keys to return.
   * @returns A promise that resolves to an array of issue keys that match the search criteria.
   * @throws Will throw an error if the search operation fails.
   */
  export async function getIssueKeys(client: jira.Version3Client, projectKey: string, issueType: string, summary: string, status: string, issueKeysToIgnore: string[], maxResults: number): Promise<string[]> {
    try {
      let jql = `project = ${projectKey} AND issuetype = "${issueType}" AND summary ~ "${summary}" AND status = "${status}"`;
      if (issueKeysToIgnore.length > 0) {
        jql = `${jql} AND issuekey NOT IN (${issueKeysToIgnore.join(',')})`
      }
      core.debug(`Searching for issue using jql: '${jql}'`)
      const result = await client.issueSearch.searchForIssuesUsingJql({
        jql: jql,
        maxResults: maxResults,
        fields: ['key']
      });

      if (result.issues === undefined) {
        core.debug('Found no matching issues.')
        return [];
      }

      return result.issues.map(issue => issue.key);
    } catch (error) {
      core.error(`Error searching for issues: ${error}`);
      throw error
    }
  }

  /**
   * Links two Jira issues with a specified link type.
   *
   * This function creates a link between an inward issue and an outward issue using the provided link type.
   * It logs the linking process and handles any errors that may occur during the operation.
   *
   * @param client - The Jira client instance used to perform the linking.
   * @param inwardIssueKey - The key of the inward issue to be linked.
   * @param outwardIssueKey - The key of the outward issue to be linked.
   * @param linkType - The type of link to create between the issues (e.g., 'Blocks', 'Relates').
   * @throws Will throw an error if the linking operation fails.
   */
  export async function linkIssues(client: jira.Version3Client, inwardIssueKey: string, outwardIssueKey: string, linkType: string) {
    core.debug(`Linking issue ${inwardIssueKey} to ${outwardIssueKey} with link type of '${linkType}'`)
    try {
      await client.issueLinks.linkIssues({
        type: {
          name: linkType,
        },
        inwardIssue: {
          key: inwardIssueKey,
        },
        outwardIssue: {
          key: outwardIssueKey,
        }
      });

      core.debug(`Successfully linked issues: ${inwardIssueKey} now '${linkType}' ${outwardIssueKey}`);
    } catch (error) {
        core.error(`Error linking issues ${inwardIssueKey} and ${outwardIssueKey}: ` + error);
        throw error
    }
  }

  /**
   * Creates a new Solidity Review issue in the specified Jira project.
   *
   * This function clones an existing issue to create a new Solidity Review issue in the given project.
   * It also clones all linked issues and their checklists. If any error occurs during the process,
   * it attempts to clean up any partially created issues.
   *
   * @param client - The Jira client instance used to perform the operations.
   * @param projectKey - The key of the project where the new issue will be created.
   * @param sourceIssueKey - The key of the issue to be cloned as the new Solidity Review issue.
   * @returns A promise that resolves to the key of the newly created Solidity Review issue.
   * @throws Will throw an error if the creation or cloning process fails.
   */
  export async function createSolidityReviewIssue(client: jira.Version3Client, projectKey: string, sourceIssueKey: string) {
    let solidityReviewKey = ""
    try {
      core.info(`Creating new Solidity Review issue in project ${projectKey}`)
      solidityReviewKey = await cloneIssue(client, sourceIssueKey, projectKey)
      await cloneLinkedIssues(client, projectKey, sourceIssueKey, solidityReviewKey, ['Epic'], 2, 1)
      core.info(`Created new Solidity Review issue in project ${projectKey}. Issue key: ${solidityReviewKey}`)

      return solidityReviewKey
    } catch (error) {
      core.setFailed('Failed to create new Solidity Review issue: ' + error)
      if (solidityReviewKey !== '') {
        await cleanUpUnfinishedIssues(client, [solidityReviewKey])
      }
      throw error
    }
  }

  /**
   * Clones an existing Jira issue into a specified project.
   *
   * This function retrieves the details of an existing issue and creates a new issue in the specified project
   * with the same details (priority, summary, description, and issue type).
   *
   * @param client - The Jira client instance used to perform the operations.
   * @param originalIssueKey - The key of the issue to be cloned.
   * @param projectKey - The key of the project where the new issue will be created.
   * @returns A promise that resolves to the key of the newly created issue.
   * @throws Will throw an error if the cloning process fails.
   */
  export async function cloneIssue(client: jira.Version3Client, originalIssueKey: string, projectKey: string): Promise<string> {
    try {
      core.debug(`Trying to clone ${originalIssueKey}`)
      const originalIssue = await client.issues.getIssue({ issueIdOrKey: originalIssueKey });

      if (originalIssue.fields.issuetype === undefined) {
        throw new Error(`Issue ${originalIssueKey} is missing issue type id. This should not happen.`)
      }

      const newIssue = await client.issues.createIssue({
        fields: {
          project: {
            key: projectKey,
          },
          priority: originalIssue.fields.priority,
          summary: originalIssue.fields.summary,
          description: originalIssue.fields.description,
          issuetype: { id: originalIssue.fields.issuetype.id },
        },
      });
      core.debug(`Cloned issue key: ${newIssue.key}`)
      return newIssue.key;
    } catch (error) {
      core.error(`Error cloning issue ${originalIssueKey}: ` + error)
      throw error
    }
  }

  /**
   * Clones all linked issues of specified types from a source issue to a target issue within a given project.
   *
   * This function retrieves all linked issues of the specified types from the source issue, clones them into the target project,
   * copies their checklists, and links them to the target issue. If any error occurs during the process, it attempts to clean up
   * any partially created issues.
   *
   * @param client - The Jira client instance used to perform the operations.
   * @param projectKey - The key of the project where the new issues will be created.
   * @param sourceIssueKey - The key of the source issue whose linked issues will be cloned.
   * @param targetIssueKey - The key of the target issue to which the cloned issues will be linked.
   * @param issueTypes - An array of issue types to filter the linked issues (e.g., 'Task', 'Bug').
   * @param expectedLinkedIssues - The expected number of linked issues to be cloned.
   * @param expectedMinChecklists - The minimum number of checklists each issue should have.
   * @throws Will throw an error if the cloning process fails or if the number of linked issues does not match the expected count.
   */
  export async function cloneLinkedIssues(client: jira.Version3Client, projectKey: string, sourceIssueKey: string, targetIssueKey: string, issueTypes: string[], expectedLinkedIssues: number, expectedMinChecklists: number) {
    const linkedIssuesKeys: string[] = []
    try {
        core.debug(`Cloning to ${targetIssueKey} all issues with type '${join(...issueTypes)}' linked to ${sourceIssueKey}`)
      const originalIssue = await client.issues.getIssue({ issueIdOrKey: sourceIssueKey });

      const linkedIssues = originalIssue.fields.issuelinks.filter(link => {
        const issueTypeName = link.inwardIssue?.fields?.issuetype?.name
        if (!issueTypeName) {
          return false
        }
        return issueTypes.length === 0 || issueTypes.includes(issueTypeName)
      });

      if (linkedIssues.length !== expectedLinkedIssues) {
        throw new Error(`Expected exactly ${expectedLinkedIssues} linked issues of type ${join(...issueTypes)}, but got ${linkedIssues.length}`)
      }

      for (const issueLink of linkedIssues) {
        if (!issueLink.inwardIssue?.key) {
          throw new Error(`Issue link ${issueLink.id} was missing inward issue or inward issue key`)
        }

        const linkedIssue = await client.issues.getIssue({ issueIdOrKey: issueLink.inwardIssue?.key });

        if (linkedIssue.fields.issuetype === undefined) {
            throw new Error(`Issue ${linkedIssue.key} is missing issue type id. This should not happen.`)
          }

        const newLinkedIssue = await client.issues.createIssue({
          fields: {
            project: {
              key: projectKey,
            },
            priority: linkedIssue.fields.priority,
            summary: linkedIssue.fields.summary,
            description: linkedIssue.fields.description,
            issuetype: { id: linkedIssue.fields.issuetype.id },
          },
        });
        linkedIssuesKeys.push(newLinkedIssue.key)

        core.debug(`Cloned linked issue key: ${newLinkedIssue.key}`);

        await copyAllChecklists(linkedIssue.id, newLinkedIssue.id, expectedMinChecklists)

        await client.issueLinks.linkIssues({
          type: { name: 'Blocks' },
          inwardIssue: { key: newLinkedIssue.key },
          outwardIssue: { key: targetIssueKey },
        });

        core.debug(`Linked ${newLinkedIssue.key} to issue ${targetIssueKey}`);
      }
    } catch (error) {
        core.error(`Error cloning linked issues from ${sourceIssueKey} to ${targetIssueKey}:  ${error}`);
        core.info('issues so far: ' + linkedIssuesKeys)
        await cleanUpUnfinishedIssues(client, linkedIssuesKeys)
        throw error
    }
  }

  /**
   * Cleans up unfinished Jira issues by closing them with a 'Declined' resolution.
   *
   * This function iterates over the provided issue keys and attempts to close each issue with a 'Declined' resolution.
   * If any error occurs during the process, it logs the error and returns it.
   *
   * @param client - The Jira client instance used to perform the operations.
   * @param issueKeys - An array of issue keys to be closed.
   * @returns A promise that resolves to `undefined` if all issues are closed successfully, or an error if any issue fails to close.
   * @throws Will throw an error if closing any of the issues fails.
   */
  async function cleanUpUnfinishedIssues(client: jira.Version3Client, issueKeys: string[]): Promise<unknown> {
    try {
      for (const key of issueKeys) {
        await declineIssue(client, key, 'Closing issue due to an error in automatic creation of Solidity Review')
      }
      return
    } catch (error) {
      core.error(`Failed to close at least one of issues: ${join(...issueKeys)} due to: ${error}. Please close them manually`)
      return error
    }
  }

  /**
   * Declines a Jira issue by transitioning it to a 'Closed' status with a 'Declined' resolution and adding a comment.
   *
   * This function transitions the specified Jira issue to the 'Closed' status using the provided transition ID and resolution.
   * It also adds a comment to the issue explaining the reason for the decline.
   *
   * @param client - The Jira client instance used to perform the operations.
   * @param issueKey - The key of the issue to be declined.
   * @param commentText - The text of the comment to be added to the issue.
   * @throws Will throw an error if the transition or comment operation fails.
   */
  async function declineIssue(client: jira.Version3Client, issueKey: string, commentText: string) {
    // in our JIRA '81' is transitionId of `Closed` status, using transition name did not work
    await transitionIssueWithComment(client,issueKey, '81', 'Declined', commentText)
  }

  /**
   * Transitions a Jira issue to a specified status with a given resolution and adds a comment.
   *
   * This function transitions the specified Jira issue to the status identified by the provided transition ID.
   * It also sets the resolution of the issue and adds a comment explaining the reason for the transition.
   *
   * @param client - The Jira client instance used to perform the operations.
   * @param issueKey - The key of the issue to be transitioned.
   * @param transitionId - The ID of the transition to be applied to the issue.
   * @param resolution - The resolution to be set for the issue (e.g., 'Fixed', 'Declined').
   * @param commentText - The text of the comment to be added to the issue.
   * @throws Will throw an error if the transition or comment operation fails.
   */
  export async function transitionIssueWithComment(client: jira.Version3Client, issueKey: string, transitionId: string, resolution: string, commentText: string) {
    try {
      await client.issues.doTransition({
        issueIdOrKey: issueKey,
        transition: {
          id: transitionId
        },
        fields: {
          resolution: {
            name: resolution
          }
        },
        update: {
          comment: [
            {
              add: {
                body: {
                  type: 'doc',
                  version: 1,
                  content: [
                    {
                      type: 'paragraph',
                      content: [
                        {
                          type: 'text',
                          text: commentText
                        }
                      ]
                    }
                  ]
                }
              }
            }
          ]
        }
      });

      core.debug(`Issue ${issueKey} successfully closed with comment.`);
    } catch (error) {
      core.error(`Failed to update issue ${issueKey}: ${error}`);
      throw error
    }
  }

  /**
   * Copies all checklists from a source Jira issue to a target Jira issue.
   *
   * This function retrieves the checklists from the source issue, verifies that the number of checklists meets the expected minimum,
   * and then adds the checklists to the target issue.
   *
   * @param sourceIssueId - The ID of the source issue from which checklists will be copied.
   * @param targetIssueId - The ID of the target issue to which checklists will be added.
   * @param expectedMinChecklists - The minimum number of checklists expected in the source issue.
   * @throws Will throw an error if the number of checklists in the source issue is less than the expected minimum or if any operation fails.
   */
  export async function copyAllChecklists(sourceIssueId: string, targetIssueId: string, expectedMinChecklists: number) {
    core.debug(`Copying all checklists from ${sourceIssueId} to ${targetIssueId}`)
    const checklistProperty = 'sd-checklists-0'
    const checklistJson = await getChecklistJSONFromIssue(sourceIssueId, checklistProperty)
    assertChecklistCount(checklistJson, expectedMinChecklists)
    addChecklistsToIssue(targetIssueId, checklistProperty, checklistJson)
    core.debug(`Copied all checklists from ${sourceIssueId} to ${targetIssueId}`)
  }

  /**
   * Asserts whether there are at least a specified number of checklists in the provided checklist JSON.
   *
   * This function checks if the `checklists` array in the provided JSON contains at least the specified minimum number of checklists.
   * If the array is missing or contains fewer checklists than expected, it throws an error.
   *
   * Sample checklist:
   * {
       "version":"1.0.0",
       "checklists":[
          {
            "id":0,
            "name":"Example to do",
            "items":[
                {
                  "name":"<p>task 1</p>",
                  "required":true,
                  "completed":true,
                  "status":0,
                  "user":"{ user id }",
                  "date":"2019-08-13T13:18:24.046Z"
                },
                {
                  "name":"<p>task 2</p>",
                  "required":false,
                  "completed":true,
                  "status":0,
                  "user":"{ user id }",
                  "date":"2019-08-13T13:18:44.988Z"
                },
                {
                  "name":"<p>task 3</p>",
                  "required":false,
                  "completed":false,
                  "status":0,
                  "user":"{ user id }",
                  "date":"2019-08-08T13:30:29.643Z"
                }
            ]
          }
       ]
    }
   * @param checklistJSON - The JSON object containing the checklists to be validated.
   * @param minChecklistCount - The minimum number of checklists expected in the `checklists` array.
   * @throws Will throw an error if the `checklists` array is missing, not an array, or contains fewer checklists than `minChecklistCount`.
   */
  export function assertChecklistCount(checklistJSON: any, minChecklistCount: number) {
    if (checklistJSON.checklists) {
      if (!(checklistJSON.checklists instanceof Array) || (checklistJSON.checklists as Array<any>).length < minChecklistCount) {
        core.debug('Checklist JSON:')
        core.debug(JSON.stringify(checklistJSON))
        throw new Error(`Checklist JSON either did not contain "checklists" array or it's lenght was smaller than ${minChecklistCount}`)
      }

      return
    }

    core.debug('Checklist JSON:')
    core.debug(JSON.stringify(checklistJSON))
    throw new Error('Checklist JSON did not contain any checklists.')
  }

  /**
   * Adds provided checklists to a Jira issue. Works only with checklists created with the plugin "Multiple checklists for Jira".
   * It's designed to work with the output of the `getChecklistJSONFromIssue()` function.
   *
   * This function sends a PUT request to the Jira API to add the checklists to the specified issue.
   * It logs the process and handles any errors that may occur during the operation.
   *
   * @param issueId - The Jira issue ID (not key) of the issue to which checklists will be added.
   * @param checklistProperty - The name of the checklist property, usually in the format `sd-checklists-{N}`.
   * @param checklistsJson - The JSON object containing the checklists to be added, conforming to the "Multiple checklists for Jira" format.
   * @throws Will throw an error if the operation to add checklists fails.
   */
  export async function addChecklistsToIssue(issueId: string, checklistProperty: string, checklistsJson: object) {
    core.debug(`Adding checklists to issue ${issueId}`)
    const { jiraHost, jiraUserName, jiraApiToken } = getJiraEnvVars();

    try {
    await axios.put(
        `${jiraHost}rest/api/3/issue/${issueId}/properties/${checklistProperty}`,
        checklistsJson,
        {
        auth: {
            username: jiraUserName,
            password: jiraApiToken,
        }
    },
    );
    core.debug(`Added checklists successfully`)
    } catch (error) {
        core.error(`Failed to add checklists to issue ${issueId}: ${error}`)
        throw error
    }
  }

  /**
   * Reads all checklists created with the plugin "Multiple checklists for Jira" and returns them as JSON.
   * This JSON can be used as-is to add these checklists to another issue. It's meant to be used in tandem
   * with `addChecklistsToIssue()`.
   *
   * This function sends a GET request to the Jira API to fetch the checklists from the specified issue.
   * It logs the process and handles any errors that may occur during the operation.
   *
   * @param issueId - The Jira issue ID (not key) of the issue to check.
   * @param checklistProperty - The name of the checklist property, usually in the format `sd-checklists-{N}`.
   * @returns A promise that resolves to a JSON object containing the checklists.
   * @throws Will throw an error if the operation to fetch checklists fails or if the response has unexpected content.
   */
  export async function getChecklistJSONFromIssue(issueId: string, checklistProperty: string): Promise<object> {
    core.debug(`Fetching all checklists from issue ${issueId}`)
    const { jiraHost, jiraUserName, jiraApiToken } = getJiraEnvVars();

    try {
        const response = await axios.get(
          `${jiraHost}rest/api/3/issue/${issueId}/properties/${checklistProperty}`,
          {
            auth: {
              username: jiraUserName,
              password: jiraApiToken,
            },
          }
        );

        if (response.data.value?.checklists && response.data.value?.checklists instanceof Array) {
            core.debug(`Found ${(response.data.value?.checklists as Array<unknown>).length} checklists`)

            return response.data.value
        }

        throw new Error('Checklist response had unexpected content: ' + JSON.stringify(response.data))

      } catch (error) {
        core.error(`Error reading checklists from issue ${issueId}: ${error}`);

        throw error
      }
  }

  /**
   * Queries Jira for open Solidity Review issues within a specified project.
   *
   * This function searches for issues within the given project that match the 'Initiative' issue type,
   * contain 'Solidity Review' in their summary, and have an 'Open' status. It excludes any issue keys
   * provided in the `issueKeysToIgnore` array and limits the number of results to 10.
   *
   * @param client - The Jira client instance used to perform the search.
   * @param projectKey - The key of the project to search within.
   * @param issueKeysToIgnore - An array of issue keys to exclude from the search results.
   * @returns A promise that resolves to an array of issue keys that match the search criteria.
   * @throws Will throw an error if the search operation fails.
   */
  async function getOpenSolidityReviewIssuesForProject(client: jira.Version3Client, projectKey: string, issueKeysToIgnore: string[]): Promise<string[]> {
    //TODO: change 'Initiative' to 'Solidity Review' once it has been created
    const issueKeys = await getIssueKeys(client, projectKey, 'Initiative', 'Solidity Review', 'Open', issueKeysToIgnore, 10)
    core.info(`Found ${issueKeys.length} open Solidity Review issues for project ${projectKey}`)
    return issueKeys
  }

  function extracProjectFromIssueKey(issueKey: string): string | undefined {
    const pattern = /([A-Z]{2,})-\d+/

    const match = issueKey.toUpperCase().match(pattern);
    const projectExtracted = match ? match[1] : undefined

    core.debug(`Extracted following project '${projectExtracted}' from issue '${issueKey}'`)

    return projectExtracted
  }

  /**
   * Reads Jira issue key of Solidity Review blueprint.
   *
   * This function retrieves the issue key from the environment variable `SOLIDITY_REVIEW_TEMPLATE_KEY`.
   * The issue key is used as a blueprint for creating new Solidity Review issues.
   *
   * @returns {string} The key of the issue that will be used as a blueprint for creating new Solidity Review issues.
   * @throws {Error} If the `SOLIDITY_REVIEW_TEMPLATE_KEY` environment variable is not set or is empty.
   */
  function readSolidityReviewTemplateKey(): string {
    const issueKey = process.env.SOLIDITY_REVIEW_TEMPLATE_KEY;
    if (!issueKey) {
      throw Error("Missing required environment variable SOLIDITY_REVIEW_TEMPLATE_KEY");
    }

    return issueKey
  }

  /**
   * Exports Jira issue keys to GitHub environment variables if the `EXPORT_JIRA_ISSUE_KEYS` environment variable is set to 'true'.
   *
   * This function checks if the `EXPORT_JIRA_ISSUE_KEYS` environment variable is set to 'true'. If it is, it exports the provided
   * Jira issue keys to the GitHub environment by appending them to the `GITHUB_ENV` file. The environment variables used are
   * `PR_JIRA_ISSUE_KEY` for the pull request issue key and `SOLIDITY_REVIEW_ISSUE_KEY` for the Solidity review issue key.
   *
   * @param prIssueKey - The Jira issue key representing the pull request.
   * @param solidityReviewIssueKey - The Jira issue key representing the Solidity review.
   */
  function exportIssueKeysToGithubEnv(prIssueKey: string, solidityReviewIssueKey: string) {
    const shouldExport = process.env.EXPORT_JIRA_ISSUE_KEYS;
    if (!shouldExport || shouldExport !== 'true' ) {
      return
    }

    const prIssueKeyEnvVar = 'PR_JIRA_ISSUE_KEY'
    const solidityReviewIssueKeyEnvVar = 'SOLIDITY_REVIEW_JIRA_ISSUE_KEY'

    fs.appendFileSync(process.env.GITHUB_ENV ?? '', `${prIssueKeyEnvVar}=${prIssueKey}\n`);
    fs.appendFileSync(process.env.GITHUB_ENV ?? '', `${solidityReviewIssueKeyEnvVar}=${solidityReviewIssueKey}\n`);

    core.info(`Exported Jira issue key representing the PR as ${prIssueKeyEnvVar} env var`)
    core.info(`Exported Jira issue key representing the Solidity Review as ${solidityReviewIssueKeyEnvVar} env var`)
  }

  async function run() {
    try {
      await main();
    } catch (error) {
      if (error instanceof Error) {
        return core.setFailed(error.message);
      }
      core.setFailed(error as any);
    }
  }

  run();
