import expect from 'expect-puppeteer'
import { Server } from 'http'
import { Browser, launch, Page } from 'puppeteer'
import { closeDbConnection, getDb } from '../src/database'
import { createChainlinkNode } from '../src/entity/ChainlinkNode'
import { JobRun } from '../src/entity/JobRun'
import { DEFAULT_TEST_PORT, start as startServer } from '../src/support/server'
import { createJobRun } from '../src/factories'

describe('End to end', () => {
  let browser: Browser
  let page: Page
  let server: Server

  beforeAll(async () => {
    browser = await launch({
      args: ['--no-sandbox'],
      devtools: false,
      headless: true,
    })

    page = await browser.newPage()
    server = await startServer()

    page.on('console', msg => console.log('PAGE LOG:', msg.text()))
  })

  afterAll(async () => {
    return Promise.all([browser.close(), server.close(), closeDbConnection()])
  })

  it('can search for job run', async () => {
    const db = await getDb()
    const [node, _] = await createChainlinkNode(db, 'endToEndChainlinkNode')
    const jobRun = await createJobRun(db, node)

    await page.goto(`http://localhost:${DEFAULT_TEST_PORT}`)
    await expect(page).toFill('form input[name=search]', jobRun.runId)
    await expect(page).toClick('form button')
    await page.waitForNavigation()
    await expect(page).toMatch(jobRun.runId)
  })
})
