/* eslint-disable @typescript-eslint/no-var-requires */

const pupExpect = require('expect-puppeteer')
const PupHelper = require('../support/PupHelper.js')

describe('End to end', () => {
  let browser, page, pupHelper

  beforeAll(async () => {
    ;({ browser, page, pupHelper } = await PupHelper.launch())
  })

  afterAll(async () => {
    return browser.close()
  })

  it('gets the configuration page', async () => {
    await pupHelper.signIn()

    // Visit config page
    await pupHelper.clickLink('Configuration')
    await pupExpect(page).toMatchElement('h5', { text: 'Configuration' })
    await pupExpect(page).toMatch('ACCOUNT_ADDRESS')
  })
})
