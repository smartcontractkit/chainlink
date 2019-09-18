/* eslint-disable @typescript-eslint/no-var-requires */

const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const { newServer } = require('../support/server.js')
const { scrape } = require('../support/scrape.js')
const puppeteerConfig = require('../puppeteer.config.js')
const PupHelper = require('../support/PupHelper.js')
const generateJobJson = require('../support/generateJobJson.js')

describe('End to end', () => {
  let browser, page, server, pupHelper
  beforeAll(async () => {
    jest.setTimeout(30000)
    pupExpect.setDefaultOptions({ timeout: 3000 })
    server = await newServer(`{"last": "3843.95"}`)
    browser = await puppeteer.launch(puppeteerConfig)
    page = await browser.newPage()
    pupHelper = new PupHelper(page)
  })

  afterAll(async () => {
    return Promise.all([browser.close(), server.close()])
  })

  it('duplicates a job', async () => {
    await pupHelper.signIn()

    // Create Job
    await pupHelper.clickLink('New Job')
    await pupHelper.waitForContent('h5', 'New Job')

    const jobJson = generateJobJson(server.port)
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
