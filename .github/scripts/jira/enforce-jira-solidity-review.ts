import * as core from "@actions/core";
import jira from "jira.js";
import axios, { AxiosError, isAxiosError } from "axios";
import { join } from "path";
import { createJiraClient, extractJiraIssueNumbersFrom, getJiraEnvVars, PR_PREFIX, SOLIDITY_REVIEW_PREFIX } from "./lib";
import { appendIssueNumberToChangesetFile, extractChangesetFile } from "./changeset-lib";

const SHARED_PRODUCT = 'shared'
const SOLIDITY_REVIEW_TEMAPLTE_KEY = 'TT-1634'

async function main() {
    // const dryRun = !!process.env.DRY_RUN;
    const { productToProjectMap, affectedProducts } = fetchEnvironmentVariables();
    const { changesetFile } = extractChangesetFile();

    // first let's make sure that this PR is linked to a JIRA issue, we will need it later anyway
    const jiraPRIssues = await extractJiraIssueNumbersFrom(PR_PREFIX, [changesetFile])
    if (jiraPRIssues.length !== 1) {
        core.setFailed(
            `This function can only work with 1 JIRA issue related to PRs, but found ${jiraPRIssues.length}`
          );
          return false
    }

    const jiraPRIssue = jiraPRIssues[0]

    // now let's check whether the issue is already linked to at least one Solidity Review issue (it's okay if there's more than one if PR modifies files for more than one project)
    const jiraSolidityIssues = await extractJiraIssueNumbersFrom(SOLIDITY_REVIEW_PREFIX, [changesetFile])
    if (jiraSolidityIssues.length > 1) {
        core.info(`Found linked Solidity Review issues ${join(...jiraSolidityIssues)}. Nothing more needs to be done.`)

        return false
    }

    const productToJiraProjectMap = parseProductToProjectMap(productToProjectMap)
    let products = parseAffectedProducts(affectedProducts)

    const filteredProducts = products.filter(product => product !== SHARED_PRODUCT)
    let jiraProject: string


    // if only shared contracts were modified we expect them to be reviewed within the Jira project, to which the PR belongs
    if (filteredProducts.length == 0) {
        // throw new Error(`Only ${SHARED_PRODUCT} contracts were modified, it's impossible to say what project link it to.`)
        const maybeJiraProject = extracProjectFromIssueKey(jiraPRIssue)
        if (!maybeJiraProject) {
            throw new Error(`Could not extract project key from issue: ${jiraPRIssue}`)
        }
        jiraProject = maybeJiraProject
    } else if (filteredProducts.length === 1) {
        jiraProject = productToJiraProjectMap[filteredProducts[0]]
    // PR modifies more than 1 product, that's not supported currently, but in reality it should be?
    } else {
        throw new Error(`PR should modify Solidity files only for one product, but changes to ${filteredProducts.length} were detected.`)
    }

    if (!jiraProject) {
        throw new Error(`No JIRA project found for product ${filteredProducts[0]}. Please provide missing mapping and re-run this workflow.`)
    }

    const client = createJiraClient();

    var solidityReviewIssueKey: string
    const openSolidityReviewIssues = await getOpenSolidityReviewIssuesForProject(client, jiraProject, [SOLIDITY_REVIEW_TEMAPLTE_KEY])
    if (openSolidityReviewIssues.length == 0) {
        solidityReviewIssueKey = await createSolidityReviewIssue(client, jiraProject, SOLIDITY_REVIEW_TEMAPLTE_KEY)
    } else if (openSolidityReviewIssues.length == 1) {
        solidityReviewIssueKey = openSolidityReviewIssues[0]
    } else {
        throw new Error(`Found following open Solidity Review issues for project ${jiraProject}: ${join(...openSolidityReviewIssues)}. This is incorrect, there should ever only be 1 open issue of this type`)
    }

    // this should just throw error, instead of returning a bool
    const isBlocksLinkSet = linkIssues(client, solidityReviewIssueKey, jiraPRIssue, 'Blocks')
    if (!isBlocksLinkSet) {
        throw new Error(`Failed to block issue ${jiraPRIssue} by ${solidityReviewIssueKey}`)
    }

    core.info(`Appending JIRA Solidity Review issue ${solidityReviewIssueKey} to changeset file`);
    await appendIssueNumberToChangesetFile(SOLIDITY_REVIEW_PREFIX, changesetFile, solidityReviewIssueKey);
  }

  function parseProductToProjectMap(input: string): Record<string, string> {
    const productToJiraProjectMap: Record<string, string> = {};

    input.split(",").forEach(pair => {
        const [product, jiraProject] = pair.split('=')
        productToJiraProjectMap[product] = jiraProject
    })

    return productToJiraProjectMap
  }

  function parseAffectedProducts(input: string): string[] {
    return input.split(',')
  }

  async function getOpenSolidityReviewIssuesForProject(client: jira.Version3Client, projectKey: string, issueKeysToIgnore: string[]): Promise<string[]> {
    //TODO: change 'Initiative' to 'Solidity Review' once it has been created
    return getIssueKeys(client, projectKey, 'Initiative', 'Solidity Review', 'Open', issueKeysToIgnore, 10)
  }

  function extracProjectFromIssueKey(issueKey: string): string | undefined {
    const pattern = /([A-Z]{2,})-\d+/

    const match = issueKey.toUpperCase().match(pattern);
    return match ? match[1] : undefined;
  }

  async function getIssueKeys(client: jira.Version3Client, projectKey: string, issueType: string, title: string, status: string, issueKeysToIgnore: string[], maxResults?: number): Promise<string[]> {
    if (maxResults === undefined) {
        maxResults = 10
    }
    try {
      let jql = `project = ${projectKey} AND issuetype = "${issueType}" AND summary ~ "${title}" AND status = "${status}"`;
      if (issueKeysToIgnore.length > 0) {
        jql = `${jql} AND issuekey NOT IN (${issueKeysToIgnore.join(',')})`
      }
      const result = await client.issueSearch.searchForIssuesUsingJql({
        jql: jql,
        maxResults: maxResults,
        fields: ['key']
      });

      if (result.issues == undefined) {
        core.info('Found no matching issues.')
        return [];
      }

      return result.issues.map(issue => issue.key);
    } catch (error) {
      core.error('Error searching for issue:', error);
      return [];
    }
  }

  export async function linkIssues(client: jira.Version3Client, inwardIssueKey: string, outwardIssueKey: string, linkType: string) {
    try {
      const response = await client.issueLinks.linkIssues({
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

      core.info(`Successfully linked issues: ${inwardIssueKey} is ${linkType} by ${outwardIssueKey}`);
      return true
    } catch (error) {
        core.error('Error linking issues:', error);
    }
    return false
  }

  async function createSolidityReviewIssue(client: jira.Version3Client, projectKey: string, sourceIssueKey: string) {
   const newIssueKey = await cloneParentIssue(client, sourceIssueKey, projectKey)
   if (!newIssueKey) {
    throw new Error('Failed cloning Solidity Review issue')
   }

    await cloneLinkedEpics(client, projectKey, sourceIssueKey, newIssueKey)

    return newIssueKey
  }

  async function cloneParentIssue(client: jira.Version3Client, originalIssueKey: string, projectKey: string): Promise<string | undefined> {
    try {
      const originalIssue = await client.issues.getIssue({ issueIdOrKey: originalIssueKey });

      if (originalIssue.fields.issuetype == undefined) {
        throw new Error(`Issue ${originalIssueKey} is missing issue type id. This should not happen.`)
      }

      const newIssue = await client.issues.createIssue({
        fields: {
          project: {
            key: projectKey,  // Target project
          },
          summary: originalIssue.fields.summary,
          description: originalIssue.fields.description,
          issuetype: { id: originalIssue.fields.issuetype.id },
          //TODO check out where checklists are stored
          customfield_12345: originalIssue.fields.customfield_12345,
        },
      });

      core.info(`Cloned issue: ${newIssue.key}`);
      return newIssue.key;
    } catch (error) {
      core.error('Error cloning parent issue:', error);
      return undefined;
    }
  }

  async function cloneLinkedEpics(client: jira.Version3Client, projectKey: string, originalIssueKey: string, newIssueKey: string) {
    try {
      const originalIssue = await client.issues.getIssue({ issueIdOrKey: originalIssueKey });

      // Check the issue's links for any linked Epics
      const linkedEpics = originalIssue.fields.issuelinks.filter(link => {
        return link.inwardIssue && link.inwardIssue.fields && link.inwardIssue.fields.issuetype && link.inwardIssue.fields.issuetype.name === 'Epic'; // Filter for Epics
      });

      if (linkedEpics.length !== 2) {
        throw new Error(`Expected exactly 2 linked epics, but got ${linkedEpics.length}`)
      }

      for (const epicLink of linkedEpics) {
        const originalEpicKey = epicLink.inwardIssue!!.key!!;

        // Clone each linked Epic
        const originalEpic = await client.issues.getIssue({ issueIdOrKey: originalEpicKey });

        if (originalEpic.fields.issuetype == undefined) {
            throw new Error(`Issue ${originalEpic} is missing issue type id. This should not happen.`)
          }

        const newEpic = await client.issues.createIssue({
          fields: {
            project: {
              key: projectKey,
            },
            priority: originalEpic.fields.priority,
            summary: originalEpic.fields.summary,
            description: originalEpic.fields.description,
            issuetype: { id: originalEpic.fields.issuetype.id },
          },
        });

        core.info(`Linked Epic cloned: ${newEpic.key}`);

        copyAllChecklists(originalEpic.id, newEpic.id)

        await client.issueLinks.linkIssues({
          type: { name: 'Blocks' },
          inwardIssue: { key: newEpic.key },
          outwardIssue: { key: newIssueKey },
        });

        core.info(`Linked Epic ${newEpic.key} to custom issue ${newIssueKey}`);
      }
    } catch (error) {
        core.error('Error cloning linked epics:', error);
        throw error
    }
  }

  async function copyAllChecklists(sourceIssueId: string, targetIssueId: string): Promise<Error | undefined> {
    try {
        const checklistProperty = 'sd-checklists-0'
        const checklistJson = await getChecklistsFromIssue(sourceIssueId, checklistProperty)
        addChecklistsToIssue(targetIssueId, checklistProperty, checklistJson)
    } catch(error) {
        // it means that there are no more checklists to copy
        if (isAxiosError(error) && error.status?.toString() == '404') {
            return undefined
        }

        return error
    }
}

  async function addChecklistsToIssue(issueId: string, checklistProperty: string, checklistsJson: JSON) {
    const { jiraHost, jiraUserName, jiraApiToken } = getJiraEnvVars();

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
  }

  async function getChecklistsFromIssue(issueId: string, checklistProperty: string): Promise<JSON> {
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

        if (response.data.value?.checklists && (response.data.value?.checklists as Array<JSON>).length > 0) {
            return response.data.value as JSON
        }

        throw new Error('Checklist response had unexpected content: ' + response.data)

      } catch (error) {
        console.error(`Error reading checklists from issueId ${issueId}:`, error);
        throw error
      }
  }

  function fetchEnvironmentVariables() {
    const productToProjectMap = process.env.PRODUCT_TO_PROJECT_MAP;
    if (!productToProjectMap) {
      throw Error("PRODUCT_TO_PROJECT_MAP environment variable is missing");
    }
    const affectedProducts = process.env.AFFECTED_PRODUCTS;
    if (!affectedProducts) {
      throw Error("AFFECTED_PRODUCTS environment variable is missing");
    }
    return { productToProjectMap, affectedProducts };
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
// async function testClone() {
//     const client = createJiraClient();
//     // await cloneLinkedEpics(client, 'TT', SOLIDITY_REVIEW_TEMAPLTE_KEY, 'TT-1637')
//     try {
//         // const originalIssue = await client.issues.getIssue({ issueIdOrKey: 'TT-1635' });

//         // console.log(originalIssue.id)
//         const r = await copyAllChecklists(425622, 425755)
//         console.log(r)
//     } catch(error) {
//         console.error(error)
//     }
// }

// testClone();