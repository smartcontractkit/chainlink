/* eslint-disable @typescript-eslint/no-var-requires */

const puppeteer = require('puppeteer')
const pupExpect = require('expect-puppeteer')
const { newServer } = require('../support/server.js')
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

  it('archives a job', async () => {
    await pupHelper.signIn()

    // Create Job
    await pupHelper.clickLink('New Job')
    await pupHelper.waitForContent('h5', 'New Job')

    const jobJson = generateJobJson(server.port)
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
