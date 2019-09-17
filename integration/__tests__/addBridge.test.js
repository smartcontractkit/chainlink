/* eslint-disable @typescript-eslint/no-var-requires */

const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const puppeteerConfig = require('../puppeteer.config.js')
const PupHelper = require('../support/PupHelper.js')

describe('End to end', () => {
  let browser, page, pupHelper
  beforeAll(async () => {
    jest.setTimeout(30000)
    pupExpect.setDefaultOptions({ timeout: 6000 })
    browser = await puppeteer.launch(puppeteerConfig)
    page = await browser.newPage()
    pupHelper = new PupHelper(page)
  })

  afterAll(async () => {
    return browser.close()
  })

  it('adds a bridge', async () => {
    // Navigate to Operator UI
    await page.goto('http://localhost:6688')
    await pupExpect(page).toMatch('Chainlink')

    // Login
    await pupHelper.signIn()
    await pupExpect(page).toMatch('Jobs')

    // Add Bridge
    // await pupHelper.clickLink('Bridges')
    await pupExpect(page).toClick('a', { text: 'Bridges' })
    await pupExpect(page).toMatchElement('h4', { text: 'Bridges' })
    // debugger
    // await pupHelper.clickLink('New Bridge')
    await pupExpect(page).toClick('a', { text: 'New Bridge' })
    await pupExpect(page).toFillForm('form', {
      name: 'new_bridge',
      url: 'http://example.com',
      minimumContractPayment: '123',
      confirmations: '5',
    })
    await pupHelper.clickButton('Create Bridge')
    await pupExpect(page).toMatch(/success.+?bridge/i)

    // Navigate to bridge show page
    const notification = await pupHelper.waitForNotification(
      'Successfully created bridge',
    )
    const notificationLink = await notification.$('a')
    await notificationLink.click()
    const pathName = await page.evaluate(() => window.location.pathname)
    expect(pathName).toEqual('/bridges/new_bridge')
  })
})
