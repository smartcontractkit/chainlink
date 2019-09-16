/* eslint-disable @typescript-eslint/no-var-requires */

const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const puppeteerConfig = require('../puppeteer.config.js')
const {
  signIn,
  consoleLogger,
  clickButton,
  clickLink,
  waitForNotification,
} = require('../support/helpers.js')

describe('End to end', () => {
  let browser, page
  beforeAll(async () => {
    jest.setTimeout(30000)
    pupExpect.setDefaultOptions({ timeout: 3000 })
    browser = await puppeteer.launch(puppeteerConfig)
    page = await browser.newPage()
    page.on('console', consoleLogger(page))
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
    await clickLink(page, 'Bridges')
    await pupExpect(page).toMatchElement('h4', { text: 'Bridges' })
    await clickLink(page, 'New Bridge')
    await pupExpect(page).toFillForm('form', {
      name: 'new_bridge',
      url: 'http://example.com',
      minimumContractPayment: '123',
      confirmations: '5',
    })
    await clickButton(page, 'Create Bridge')
    await pupExpect(page).toMatch(/success.+?bridge/i)

    // Navigate to bridge show page
    const notification = await waitForNotification(
      page,
      'Successfully created bridge',
    )
    const notificationLink = await notification.$('a')
    await notificationLink.click()
    const pathName = await page.evaluate(() => window.location.pathname)
    expect(pathName).toEqual('/bridges/new_bridge')
  })
})
