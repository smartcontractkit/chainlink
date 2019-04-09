import { Server } from 'http'
import { Browser, launch, Page } from 'puppeteer'
import expect from 'expect-puppeteer'
import { createDbConnection, closeDbConnection } from '../src/database'
import { clearDb } from '../src/__tests__/testdatabase'
import seed, { JOB_RUN_A_ID, JOB_RUN_B_ID } from '../src/seed'
import server from '../src/server'

const startServer = async () => {
  await createDbConnection()
  await seed()
  return server()
}

afterEach(async () => clearDb())

describe('End to end', () => {
  let browser: Browser, page: Page, server: Server
  beforeAll(async () => {
    browser = await launch({
      devtools: false,
      headless: true,
      args: ['--no-sandbox']
    })
    page = await browser.newPage()
    server = await startServer()
    page.on('console', msg => console.log('PAGE LOG:', msg.text()))
  })

  afterAll(async () => {
    return Promise.all([browser.close(), server.close(), closeDbConnection()])
  })

  it('can search for job run', async () => {
    await page.goto('http://localhost:8080')
    await expect(page).toFill('form input[name=search]', JOB_RUN_A_ID)
    await expect(page).toClick('form button')
    await page.waitForNavigation()
    await expect(page).toMatch(JOB_RUN_A_ID)
  })
})
