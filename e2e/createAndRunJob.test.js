const puppeteer = require('puppeteer')
const expect = require('expect-puppeteer')
const { newServer } = require('./support/server.js')
const { scrape } = require('./support/scrape.js')

describe('End to end', () => {
  let browser, page
  beforeAll(async () => {
    browser = await puppeteer.launch({
      devtools: false,
      headless: false,
      args: ['--no-sandbox']
    })
    page = await browser.newPage()
<<<<<<< HEAD
=======
    server = await newServer(`{"last": "3843.95"}`)

    page.on('console', msg => console.log('PAGE LOG:', msg.text()))
>>>>>>> Add scrape command for pulling up matches, add console
  })

  afterAll(async () => {
    await browser.close()
  })

  it('creates a job that runs', async () => {
    await page.goto('http://localhost:6688')
    await expect(page).toMatch('Chainlink')

    // Login
    await expect(page).toFill('form input[id=email]', 'notreal@fakeemail.ch')
    await expect(page).toFill('form input[id=password]', 'twochains')
    await expect(page).toClick('form button')
    await expect(page).toMatch('Jobs')

    // Create Job
    await expect(page).toClick('a', { text: 'New Job' })
    await expect(page).toMatch('New Job')

    const jobJson = `{
      "initiators": [{"type": "web"}],
      "tasks": [{"type": "NoOp"}]
    }`
    await expect(page).toFill('form textarea', jobJson)
    await expect(page).toClick('button', { text: 'Create Job' })
    await expect(page).toMatch(/success.+job/i)

    // Run Job
    await expect(page).toClick('#created-job')
    await expect(page).toMatch('Job Spec Detail')
    await expect(page).toClick('button', { text: 'Run' })
<<<<<<< HEAD
    await expect(page).toMatch(/success.+run/i)
=======
    await expect(page).toMatch(/success.+?run/i)

    // Transaction ID should eventually be coded on page like so:
    //    {"result":"0x6736ad06da823692cc66c5a51032c4aed83bfca9778eb1a7ad24de67f3f472fc"}
    const match = await scrape(page, /"result":"(0x[0-9a-f]{64})"/)
    const txHash = match[1]

    // Navigate to transactions page
    await expect(page).toClick('button', { 'aria-label': 'open drawer' })
    await expect(page).toClick('span', { text: 'Transactions' })
    await expect(page).toMatchElement('h4', { text: 'Transactions' })
    await expect(page).toMatchElement('p', { text: txHash })

    // Navigate to transaction page and check for the transaction
    await expect(page).toClick('a', { text: txHash })
    await scrape(page, txHash)
>>>>>>> Add scrape command for pulling up matches, add console
  })
})
