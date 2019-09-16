/* eslint-disable @typescript-eslint/no-var-requires */

const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const { newServer } = require('../support/server.js')
const { scrape } = require('../support/scrape.js')
const puppeteerConfig = require('../puppeteer.config.js')
const PupHelper = require('../support/PupHelper.js')

describe('End to end', () => {
  let browser, page, server, pupHelper
  beforeAll(async () => {
    jest.setTimeout(3000000)
    pupExpect.setDefaultOptions({ timeout: 3000 })
    server = await newServer(`{"last": "3843.95"}`)
    browser = await puppeteer.launch(puppeteerConfig)
    page = await browser.newPage()
    pupHelper = new PupHelper(page)
  })

  afterAll(async () => {
    return Promise.all([browser.close(), server.close()])
  })

  it('creates a job that runs', async () => {
    await page.goto('http://localhost:6688')
    await pupExpect(page).toMatch('Chainlink')

    // Login
    await pupHelper.signIn()
    await pupExpect(page).toMatch('Jobs')

    // Create Job
    await pupHelper.clickLink('New Job')
    await pupExpect(page).toMatchElement('h5', { text: 'New Job' })

    // prettier-ignore
    const jobJson = `{
      "initiators": [{"type": "web"}],
      "tasks": [
        {"type": "httpget", "params": {"get": "http://localhost:${server.port}"}},
        {"type": "jsonparse", "params": {"path": ["last"]}},
        {
          "type": "ethtx",
          "confirmations": 0,
          "params": {
            "address": "0xaa664fa2fdc390c662de1dbacf1218ac6e066ae6",
            "functionSelector": "setBytes(bytes32,bytes)"
          }
        }
      ]
    }`

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
    await pupExpect(page).toClick('a', { text: 'JSON' })
    const txHash = await scrape(page, /"transactionHash": "0x([0-9a-f]{64})"/)
    expect(txHash).toBeDefined()

    // Navigate to transactions page
    await pupHelper.clickLink('Transactions')
    await pupExpect(page).toMatchElement('h4', { text: 'Transactions' })
    await pupExpect(page).toMatchElement('p', { text: txHash })

    // Navigate to transaction page and check for the transaction
    await pupExpect(page).toClick('a', { text: txHash })

    // Navigate to Explorer
    // await new Promise(resolve => setTimeout(resolve, 5000)) // Wait for CL Node to push SyncEvent
    await page.goto('http://localhost:8080')
    await pupExpect(page).toMatch('Search')
    await pupExpect(page).toFill('form input', runId)
    await pupExpect(page).toClick('button', { text: 'Search' })

    await pupHelper.waitForContent('a', runId)
    await pupExpect(page).toMatch(runId)
    await pupExpect(page).toClick('a', { text: runId })

    await scrape(page, /Complete/)
    await pupExpect(page).toMatchElement('h5', { text: 'Complete' })
  })
})
