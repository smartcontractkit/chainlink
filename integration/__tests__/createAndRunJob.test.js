/* eslint-disable @typescript-eslint/no-var-requires */

const pupExpect = require('expect-puppeteer')
const { newServer } = require('../support/server.js')
const { scrape } = require('../support/scrape.js')
const PupHelper = require('../support/PupHelper.js')
const generateJobJson = require('../support/generateJobJson.js')

describe('End to end', () => {
  let browser, page, server, pupHelper

  beforeAll(async () => {
    server = await newServer(`{"last": "3843.95"}`)
    ;({ browser, page, pupHelper } = await PupHelper.launch())
  })

  afterAll(async () => {
    return Promise.all([browser.close(), server.close()])
  })

  it('creates a job that runs', async () => {
    await pupHelper.signIn()

    // Create Job
    await pupHelper.clickLink('New Job')
    await pupHelper.waitForContent('h5', 'New Job')

    const jobJson = generateJobJson(server.port)
    await pupExpect(page).toFill('form textarea', jobJson)
    await pupExpect(page).toClick('button', { text: 'Create Job' })
    await pupExpect(page).toMatch(/success.+job/i)

    // Run Job
    await pupExpect(page).toClick('#created-job')
    await pupExpect(page).toMatch('Job Spec Detail')
    await pupExpect(page).toClick('button', { text: 'Run' })
    await pupExpect(page).toMatch(/success.+?run/i)

    const notification = await pupHelper.waitForNotification(
      'Successfully created job run',
    )
    const notificationLink = await notification.$('a')
    const runId = await notificationLink.evaluate(tag => tag.innerText)
    await notificationLink.click()

    // Transaction ID should eventually be coded on page like so:
    //    {"result":"0x6736ad06da823692cc66c5a51032c4aed83bfca9778eb1a7ad24de67f3f472fc"}
    await pupHelper.clickLink('JSON')
    const txHash = await scrape(page, /"transactionHash": "0x([0-9a-f]{64})"/)
    expect(txHash).toBeDefined()

    // Navigate to transactions page
    await pupHelper.clickLink('Transactions')
    await pupExpect(page).toMatchElement('h4', { text: 'Transactions' })
    await pupExpect(page).toMatchElement('p', { text: txHash })

    // Navigate to transaction page and check for the transaction
    await pupHelper.clickLink(txHash)

    // Navigate to Explorer
    // await new Promise(resolve => setTimeout(resolve, 5000)) // Wait for CL Node to push SyncEvent
    await page.goto('http://localhost:8080')
    await pupExpect(page).toMatch('Search')
    await pupExpect(page).toFill('form input', runId)
    await pupExpect(page).toClick('button', { text: 'Search' })

    await pupHelper.waitForContent('a', runId)
    await pupExpect(page).toMatch(runId)
    await pupHelper.clickLink(runId)

    await scrape(page, /Complete/)
    await pupExpect(page).toMatchElement('h5', { text: 'Complete' })
  })
})
