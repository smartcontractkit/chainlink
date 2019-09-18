/* eslint-disable @typescript-eslint/no-var-requires */

const pupExpect = require('expect-puppeteer')
const PupHelper = require('../support/PupHelper.js')
const generateJobJson = require('../support/generateJobJson.js')

describe('End to end', () => {
  let browser, page, pupHelper

  beforeAll(async () => {
    ;({ browser, page, pupHelper } = await PupHelper.launch())
  })

  afterAll(async () => {
    return browser.close()
  })

  it('archives a job', async () => {
    await pupHelper.signIn()

    // Create Job
    await pupHelper.clickLink('New Job')
    await pupHelper.waitForContent('h5', 'New Job')

    const jobJson = generateJobJson(1234)
    await pupExpect(page).toFill('form textarea', jobJson)
    await pupExpect(page).toClick('button', { text: 'Create Job' })
    const notification = await pupHelper.waitForNotification(
      'Successfully created job',
    )
    const notificationLink = await notification.$('a')
    const jobId = await notificationLink.evaluate(tag => tag.innerText)

    // Archive Job
    await pupExpect(page).toClick('#created-job')
    await pupExpect(page).toMatch('Job Spec Detail')
    await pupExpect(page).toClick('button', { text: 'Archive' })
    await pupExpect(page).toMatch(/Warning/i)
    await pupExpect(page).toClick('button', { text: `Archive ${jobId}` })
    await pupHelper.waitForNotification('Successfully archived job')
  })
})
