/* eslint-disable @typescript-eslint/no-var-requires */

const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const puppeteerConfig = require('../puppeteer.config.js')
const PupHelper = require('../support/PupHelper.js')

describe('End to end', () => {
  let browser, page, pupHelper
  beforeAll(async () => {
    jest.setTimeout(30000)
    pupExpect.setDefaultOptions({ timeout: 3000 })
    browser = await puppeteer.launch(puppeteerConfig)
    page = await browser.newPage()
    pupHelper = new PupHelper(page)
  })

  afterAll(async () => {
    return browser.close()
  })

  it('gets the configuration page', async () => {
    await pupHelper.signIn()

    // Visit config page
    await pupHelper.clickLink('Configuration')
    await pupExpect(page).toMatchElement('h5', { text: 'Configuration' })
    await pupExpect(page).toClick('button', { text: 'Delete Completed Jobs' })
  })
})
