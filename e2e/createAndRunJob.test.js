const puppeteer = require('puppeteer')
const expect = require('expect-puppeteer')
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
    expect.setDefaultOptions({ timeout: 3000 })

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

    page.on('console', msg => console.log('PAGE LOG:', msg.text()))
  })

  afterAll(async () => {
    return Promise.all([browser.close(), server.close()])
  })

  it('creates a job that runs', async () => {
    await page.goto('http://localhost:6688')
    await expect(page).toMatch('Chainlink')

    // Login
    await signIn(page, 'notreal@fakeemail.ch', 'twochains')
    await expect(page).toMatch('Jobs')

    // Create Job
    await clickNewJobButton(page)
    await expect(page).toMatchElement('h5', { text: 'New Job' })

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

    await expect(page).toFill('form textarea', jobJson)
    await expect(page).toClick('button', { text: 'Create Job' })
    await expect(page).toMatch(/success.+job/i)

    // Run Job
    await expect(page).toClick('#created-job')
    await expect(page).toMatch('Job Spec Detail')
    await expect(page).toClick('button', { text: 'Run' })
    await expect(page).toMatch(/success.+?run/i)

    // Transaction ID should eventually be coded on page like so:
    //    {"result":"0x6736ad06da823692cc66c5a51032c4aed83bfca9778eb1a7ad24de67f3f472fc"}
    const match = await scrape(page, /"result":"(0x[0-9a-f]{64})"/)
    const txHash = match[1]

    // Navigate to transactions page
    await clickTransactionsMenuItem(page)
    await expect(page).toMatchElement('h4', { text: 'Transactions' })
    await expect(page).toMatchElement('p', { text: txHash })

    // Navigate to transaction page and check for the transaction
    await expect(page).toClick('a', { text: txHash })
    await scrape(page, txHash)
  })
})
