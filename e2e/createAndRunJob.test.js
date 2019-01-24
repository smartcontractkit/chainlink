const puppeteer = require('puppeteer')
const expect = require('expect-puppeteer')

describe('End to end', () => {
  let browser, page
  jest.setTimeout(100000)
  beforeAll(async () => {
    browser = await puppeteer.launch({
      devtools: false,
      headless: true,
      args: ['--no-sandbox']
    })
    page = await browser.newPage()
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
    await expect(page).toMatch(/success.+run/i)
  })
})
