const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const puppeteerConfig = require('./puppeteer.config.js')
const {
  signIn,
  consoleLogger,
  clickBridgesTab,
  clickNewBridgeButton,
} = require('./support/helpers.js')

describe('End to end', () => {
  let browser, page, server
  beforeAll(async () => {
    jest.setTimeout(30000)
    pupExpect.setDefaultOptions({ timeout: 3000 })
    browser = await puppeteer.launch(puppeteerConfig)
    page = await browser.newPage()
    page.on('console', consoleLogger)
  })

  afterAll(async () => {
    return browser.close()
  })

  it('adds a bridge', async () => {
    // Navigate to Operator UI
    await page.goto('http://localhost:6688')
    await pupExpect(page).toMatch('Chainlink')

    // Login
    await signIn(page, 'notreal@fakeemail.ch', 'twochains')
    await pupExpect(page).toMatch('Jobs')

    // Add Bridge
    await clickBridgesTab(page)
    await pupExpect(page).toMatchElement('h4', { text: 'Bridges' })
    await clickNewBridgeButton(page)
    await pupExpect(page).toFillForm('form', {
      name: 'new_bridge',
      url: 'http://example.com',
      minimumContractPayment: '123',
      confirmations: '5',
    })
    await pupExpect(page).toClick('button', { text: 'Create Bridge' })

    // Navigate to bridge show page
    const flashMessage = await page.$x(
      "//p[contains(text(), 'Successfully created bridge')]",
    )
    const jobRunLink = await flashMessage[0].$('a')
    await jobRunLink.click()
    const pathName = await page.evaluate(() => location.pathname)
    expect(pathName).toEqual('/bridges/new_bridge')
  })
})
