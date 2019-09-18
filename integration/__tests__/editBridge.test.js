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

  it('edits a bridge', async () => {
    await pupHelper.signIn()

    // Add Bridge
    await pupHelper.clickLink('Bridges')
    await pupExpect(page).toMatchElement('h4', { text: 'Bridges' })
    await pupHelper.clickLink('New Bridge')
    await pupExpect(page).toFillForm('form', {
      name: 'edit_test_bridge',
      url: 'http://example.com',
      minimumContractPayment: '123',
      confirmations: '5',
    })
    await pupExpect(page).toClick('button', { text: 'Create Bridge' })
    await pupExpect(page).toMatch(/success.+?bridge/i)

    // Navigate to bridge show page
    const notification1 = await pupHelper.waitForNotification(
      'Successfully created bridge',
    )
    const notificationLink1 = await notification1.$('a')
    await notificationLink1.click()
    await pupExpect(page).toMatch('http://example.com')
    await pupExpect(page).toMatch('123')
    await pupExpect(page).toMatch('5')

    // Edit Bridge
    await pupHelper.clickLink('Edit')
    await pupExpect(page).toFillForm('form', {
      url: 'http://example2.com',
      minimumContractPayment: '321',
      confirmations: '7',
    })
    await pupExpect(page).toClick('button', { text: 'Save Bridge' })
    const notification2 = await pupHelper.waitForNotification(
      'Successfully updated',
    )
    const notificationLink2 = await notification2.$('a')
    await notificationLink2.click()
    await pupExpect(page).toMatch('http://example2.com')
    await pupExpect(page).toMatch('321')
    await pupExpect(page).toMatch('7')
  })
})
