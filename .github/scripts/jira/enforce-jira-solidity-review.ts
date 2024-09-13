import * as core from "@actions/core";
import jira from "jira.js";
import axios from "axios";
import { join } from "path";
import { createJiraClient, extractJiraIssueNumbersFrom, getJiraEnvVars, doesIssueExist, PR_PREFIX, SOLIDITY_REVIEW_PREFIX } from "./lib";
import { appendIssueNumberToChangesetFile, extractChangesetFile } from "./changeset-lib";

async function main() {
    core.info('Started linking PR to a Solidity Review issue')
    const solidityReviewTemplateKey = readSolidityReviewTemplateKey()
    const changesetFile = extractChangesetFile();

    const jiraPRIssues = await extractJiraIssueNumbersFrom(PR_PREFIX, [changesetFile])
    if (jiraPRIssues.length !== 1) {
        core.setFailed(
            `Solidity Review enforcement only works with 1 JIRA issue per PR, but found ${jiraPRIssues.length} issues in changeset file ${changesetFile}`
          );

          return
    }

    const jiraPRIssue = jiraPRIssues[0]
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

    const jiraProject = extracProjectFromIssueKey(jiraPRIssue)
    if (!jiraProject) {
        core.setFailed(`Could not extract project key from issue: ${jiraPRIssue}`)

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
    await linkIssues(client, solidityReviewIssueKey, jiraPRIssue, 'Blocks')

    core.info(`Appending JIRA Solidity Review issue ${solidityReviewIssueKey} to changeset file`);
    await appendIssueNumberToChangesetFile(SOLIDITY_REVIEW_PREFIX, changesetFile, solidityReviewIssueKey);
    core.info('Finished linking PR to a Solidity Review issue')
  }

  /**
   * Searches Jira for issues with given type and status inside a single project. Summary is isn't matched exactly, but instead must contain given string.
   * You can pass optional list of issue keys that should be excluded from search and maximum number of results.
   *
   * @param client jira client
   * @param projectKey project to search in
   * @param issueType issue type, e.g. 'Task', 'Epic'
   * @param summary summary or title of the issue
   * @param status issue status, e.g. 'Open' , 'In progress'
   * @param issueKeysToIgnore keys of issues to exclude from search
   * @param maxResults maximum number of results to return
   * @returns list of issue keys found (empty array if no match is found)
   * @throws {Error} if the search fails
   */
  async function getIssueKeys(client: jira.Version3Client, projectKey: string, issueType: string, summary: string, status: string, issueKeysToIgnore: string[], maxResults: number): Promise<string[]> {
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
   * Links two issue types with the given link type.
   *
   * For example, to block issue A by issue B. You would use link type of 'Blocks' and pass 'A' as inward issue and 'B' as outward.
   *
   * @param client jira client
   * @param inwardIssueKey
   * @param outwardIssueKey
   * @param linkType name of link to create, e.g. 'Blocks', 'Relates'
   */
  async function linkIssues(client: jira.Version3Client, inwardIssueKey: string, outwardIssueKey: string, linkType: string) {
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
   * Creates new Solidity Review issue in the given project based on the blueprint issue.
   *
   * It is also atomic in this sense, that either cloning of all issues and checklists succeds or it throws an error. If there is any
   * error during cloning it will try to close all issues that have been cloned until then, so that we don't leave Jira in undefined state.
   *
   * @param client jira client
   * @param projectKey project where to create the issue
   * @param sourceIssueKey blueprint issue
   * @returns issue key of Solidity Review created
   * @throws {Error} if creation of Solidity Review or any of its linked issues fails
   */
  async function createSolidityReviewIssue(client: jira.Version3Client, projectKey: string, sourceIssueKey: string) {
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
   * Clones Jira issue to indicated project.
   *
   * @param client jira client
   * @param originalIssueKey issue to be cloned
   * @param projectKey project to clone to
   * @returns key of the cloned issue
   * @throws {Error} if cloning fails or issue to be cloned is malformed or doesn't exist
   */
  async function cloneIssue(client: jira.Version3Client, originalIssueKey: string, projectKey: string): Promise<string> {
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
   * Clones all linked Jira between source and target issue. It can optionally limit cloning by issue types, and assert whether:
   * * source issue had at least N issues
   * * each linked issue had at least N checklists
   *
   * It's designed to be used together with "Multiple checklists for Jira" plugin and these are only checklists it supports.
   *
   * It is also atomic in this sense, that either cloning of all issues and checklists succeds or it throws an error. If there is any
   * error during cloning it will try to close all issues that have been cloned until then, so that we don't leave Jira in undefined state.
   *
   * @param client jira client
   * @param projectKey jira project key, at least two upper-cased letters
   * @param sourceIssueKey source issue key
   * @param targetIssueKey target issue key
   * @param issueTypes array of issue types to include (e.g. ['Epic', 'Task']), if empty issue type will be ignored
   * @param expectedLinkedIssues minimum number of linked issues source issue must have
   * @param expectedMinChecklists minimum number of checklists each linked issue must have
   * @throws {Error} if any of the optional expectations isn't met or if cloning fails for whatever reason
   */
  async function cloneLinkedIssues(client: jira.Version3Client, projectKey: string, sourceIssueKey: string, targetIssueKey: string, issueTypes: string[], expectedLinkedIssues: number, expectedMinChecklists: number) {
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
   * Closes as 'Declined' all issues with comment explaining that they were closed automatically during failed
   * creation of new Solidity Review issue.
   *
   * @param client jira client
   * @param issueKeys array of all issues to close
   * @returns {Error} if closing any of issues fails
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
   * Closes Jira issue with resolution equal to 'Declined'.
   *
   * @param client jira client
   * @param issueKey issue to close
   * @param commentText comment to add to the issue
   * @throws {Error} if transitioning the issue fails
   */
  async function declineIssue(client: jira.Version3Client, issueKey: string, commentText: string) {
    // in our JIRA '81' is transitionId of `Closed` status, using transition name did not work
    await transitionIssueWithComment(client,issueKey, '81', 'Declined', commentText)
  }

  /**
   * Transitions Jira issue with resolution and comment.
   *
   * @param client jria client
   * @param issueKey issue key to transit
   * @param transitionId id of transition (cannot use name!)
   * @param resolution name of resulution to use, e.x. "Won't do", "Declined", "Done", etc.
   * @param commentText comment to add to the issue
   * @throws {Error} if transitioning the issue fails
   */
  async function transitionIssueWithComment(client: jira.Version3Client, issueKey: string, transitionId: string, resolution: string, commentText: string) {
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
   * Copies all checklists between two Jira issues. It works only with "Multiple checklists for Jira" plugin.
   * It will validate whether source issue has at least N checklists.
   *
   * @param sourceIssueId jira issue key from which checklists should be copied
   * @param targetIssueId jira issue key to which checklists should be copied
   * @param expectedMinChecklists minimum number of checklists each issue should have
   * @throws {Error} if any of the issues has less checklists than expectedMinChecklists
   */
  async function copyAllChecklists(sourceIssueId: string, targetIssueId: string, expectedMinChecklists: number) {
    core.debug(`Copying all checklists from ${sourceIssueId} to ${targetIssueId}`)
    const checklistProperty = 'sd-checklists-0'
    const checklistJson = await getChecklistJSONFromIssue(sourceIssueId, checklistProperty)
    assertChecklistCount(checklistJson, expectedMinChecklists)
    addChecklistsToIssue(targetIssueId, checklistProperty, checklistJson)
    core.debug(`Copied all checklists from ${sourceIssueId} to ${targetIssueId}`)
  }

  /**
   * Asserts whether there are at least N checklists from plugin "Multiple checklists for Jira" in the checklist JSON.
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
   *
   * @param checklistJSON valid checklists array from plugin "Multiple checklists for Jira"
   * @param minChecklistCount minimum length of 'checklists' array
   * @throws {Error} if inputs field doesn't contain an Array under 'checklists' key or that array's lenght is smaller than `minChecklistCount`
   */
  function assertChecklistCount(checklistJSON: any, minChecklistCount: number) {
    if (checklistJSON.checklists) {
      if (!(checklistJSON.checklists instanceof Array) || (checklistJSON.checklists as Array<any>).length < minChecklistCount) {
        core.debug('Checklist JSON:')
        core.debug(JSON.stringify(checklistJSON))
        throw new Error(`Checklist JSON either did not contain "checklists" array or it's lenght was smaller than ${minChecklistCount}`)
      }
    }

    core.debug('Checklist JSON:')
    core.debug(JSON.stringify(checklistJSON))
    throw new Error('Checklist JSON did not contain any checklist.')
  }

  /**
   * Adds provided checklists to Jira issue. Works only with checklists created with plugin "Multiple checklists for Jira".
   * It's desgined to work with the output of `getChecklistJSONFromIssue()` function.
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
   *
   * @param issueId jira issue id (not key!) of issue to check
   * @param checklistProperty name of checklist property, usually `sd-checklists-{N}`
   * @param checklistsJson JSON of checklists conforming to "Multiple checklists for Jira" format.
   */
  async function addChecklistsToIssue(issueId: string, checklistProperty: string, checklistsJson: object) {
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
   * Reads all checklists created with plugin "Multiple checklists for Jira" and returns them as JSON
   * that can be used as-is to add these checklists to another issue. It's meant to be used in tandem
   * with `addChecklistsToIssue()`.
   *
   * @param issueId jira issue id (not key!) of issue to check
   * @param checklistProperty name of checklist property, usually `sd-checklists-{N}`
   * @returns {Promise<object>} JSON with all checklists that were found
   */
  async function getChecklistJSONFromIssue(issueId: string, checklistProperty: string): Promise<object> {
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

        throw new Error('Checklist response had unexpected content: ' + response.data)

      } catch (error) {
        core.error(`Error reading checklists from issueId ${issueId}: ${error}`);

        throw error
      }
  }

  /**
   * Queries Jira for Solidity Review issues with 'Open' status.
   *
   * @param client jira client
   * @param projectKey project symbol (at least two upper-cased letters)
   * @param issueKeysToIgnore issue keys that should be ignored during search (e.g. template blueprint)
   * @returns {Promise<string[]>} array of Solidity Review issue keys
   */
  async function getOpenSolidityReviewIssuesForProject(client: jira.Version3Client, projectKey: string, issueKeysToIgnore: string[]): Promise<string[]> {
    //TODO: change 'Initiative' to 'Solidity Review' once it has been created
    const issueKeys = await getIssueKeys(client, projectKey, 'Initiative', 'Solidity Review', 'Open', issueKeysToIgnore, 10)
    core.info(`Found ${issueKeys.length} open Solidity Review issues for project '${projectKey}'`)
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
   * @returns {string} key of issue that will be used as a blueprint for creating new Solidity Review issue.
   * @throws {Error} if SOLIDITY_REVIEW_TEMPLATE_KEY is not set or empty.
   */
  function readSolidityReviewTemplateKey(): string {
    const issueKey = process.env.SOLIDITY_REVIEW_TEMPLATE_KEY;
    if (!issueKey) {
      throw Error("Missing required environment variable SOLIDITY_REVIEW_TEMPLATE_KEY");
    }

    return issueKey
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