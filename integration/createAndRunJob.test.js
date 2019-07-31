const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const { newServer } = require('./support/server.js')
const { scrape } = require('./support/scrape.js')
const {
  signIn,
  clickNewJobButton,
  clickTransactionsMenuItem
} = require('./support/helpers.js')

const AVERAGE_CLIENT_WIDTH = 1366
const AVERAGE_CLIENT_HEIGHT = 768

describe('End to end', () => {
  let browser, page, server
  beforeAll(async () => {
    jest.setTimeout(30000)
    pupExpect.setDefaultOptions({ timeout: 3000 })
    browser = await puppeteer.launch({
      devtools: false,
      headless: true,
      args: ['--no-sandbox']
    })
    page = await browser.newPage()
    await page.setViewport({
      width: AVERAGE_CLIENT_WIDTH,
      height: AVERAGE_CLIENT_HEIGHT
    })
    server = await newServer(`{"last": "3843.95"}`)
    page.on('console', msg => {
      console.log(`PAGE LOG url: ${page.url()} | msg: ${msg.text()}`)
    })
  })

  afterAll(async () => {
    return Promise.all([browser.close(), server.close()])
  })

  it('creates a job that runs', async () => {
    await page.goto('http://localhost:6688')
    await pupExpect(page).toMatch('Chainlink')

    // Login
    await signIn(page, 'notreal@fakeemail.ch', 'twochains')
    await pupExpect(page).toMatch('Jobs')

    // Create Job
    await clickNewJobButton(page)
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
    const flashMessage = await page.$x(
      "//p[contains(text(), 'Successfully created job run')]"
    )
    await new Promise(resolve => setTimeout(resolve, 2000)) // FIXME timeout until we can reload again
    const jobRunLink = await flashMessage[0].$('a')
    const runId = await flashMessage[0].$eval('a', async link => link.innerText)
    await jobRunLink.click()

    // Transaction ID should eventually be coded on page like so:
    //    {"result":"0x6736ad06da823692cc66c5a51032c4aed83bfca9778eb1a7ad24de67f3f472fc"}
    await pupExpect(page).toClick('a', { text: 'JSON' })
    const txHash = await scrape(page, /"transactionHash": "0x([0-9a-f]{64})"/)
    expect(txHash).toBeDefined()

    // Navigate to transactions page
    await clickTransactionsMenuItem(page)
    await pupExpect(page).toMatchElement('h4', { text: 'Transactions' })
    await pupExpect(page).toMatchElement('p', { text: txHash })

    // Navigate to transaction page and check for the transaction
    await pupExpect(page).toClick('a', { text: txHash })

    // Navigate to Explorer
    await new Promise(resolve => setTimeout(resolve, 5000)) // Wait for CL Node to push SyncEvent
    await page.goto('http://localhost:8080')
    await pupExpect(page).toMatch('Search')
    await pupExpect(page).toFill('form input', runId)
    await pupExpect(page).toClick('button', { text: 'Search' })

    await new Promise(resolve => setTimeout(resolve, 500)) // FIXME not sure why we need to wait here
    await pupExpect(page).toMatch(runId)
    await pupExpect(page).toClick('a', { text: runId })

    await scrape(page, /Complete/)
    await pupExpect(page).toMatchElement('h5', { text: 'Complete' })
  })
})
