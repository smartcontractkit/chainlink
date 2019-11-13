import expect from 'expect-puppeteer'
import { Server } from 'http'
import { Browser, launch, Page } from 'puppeteer'
import { getDb } from '../src/database'
import { createChainlinkNode } from '../src/entity/ChainlinkNode'
import { DEFAULT_TEST_PORT, start, stop } from '../src/support/server'
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
    server = await start()

    page.on('console', msg => console.log('PAGE LOG:', msg.text()))
  })

  afterAll(async done => {
    browser.close()
    stop(server, done)
  })

  it('can search for job run', async () => {
    const db = await getDb()
    const [node] = await createChainlinkNode(db, 'endToEndChainlinkNode')
    const jobRun = await createJobRun(db, node)

    await page.goto(`http://localhost:${DEFAULT_TEST_PORT}`)
    await expect(page).toFill('form input[name=search]', jobRun.runId)
    await expect(page).toClick('form button')
    await page.waitForNavigation()
    await expect(page).toMatch(jobRun.runId)
  })
})
