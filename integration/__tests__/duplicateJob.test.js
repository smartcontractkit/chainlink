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

  it('duplicates a job', async () => {
    await pupHelper.signIn()

    // Create Job
    await pupHelper.clickLink('New Job')
    await pupHelper.waitForContent('h5', 'New Job')

    const jobJson = generateJobJson(1234)
    await pupExpect(page).toFill('form textarea', jobJson)
    await pupExpect(page).toClick('button', { text: 'Create Job' })
    const notification1 = await pupHelper.waitForNotification(
      'Successfully created job',
    )
    const notificationLink1 = await notification1.$('a')
    const jobId1 = await notificationLink1.evaluate(tag => tag.innerText)

    // Duplicate Job
    await pupExpect(page).toClick('#created-job')
    await pupExpect(page).toMatch('Job Spec Detail')
    await pupHelper.clickLink('Duplicate')
    await pupExpect(page).toClick('button', { text: 'Create Job' })
    const notification2 = await pupHelper.waitForNotification(
      'Successfully created job',
    )
    const notificationLink2 = await notification2.$('a')
    const jobId2 = await notificationLink2.evaluate(tag => tag.innerText)
    await pupExpect(page).toClick('#created-job')
    await pupExpect(page).toMatch('Job Spec Detail')
    expect(jobId1).not.toEqual(jobId2)
  })
})
