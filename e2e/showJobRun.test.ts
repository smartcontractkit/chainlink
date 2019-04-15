import { Server } from 'http'
import { Browser, launch, Page } from 'puppeteer'
import expect from 'expect-puppeteer'
import { closeDbConnection } from '../src/database'
import { JOB_RUN_A_ID } from '../src/seed'
import { startAndSeed as startAndSeedServer } from '../src/support/server'

describe('End to end', () => {
  let browser: Browser, page: Page, server: Server

  beforeAll(async () => {
    browser = await launch({
      devtools: false,
      headless: true,
      args: ['--no-sandbox']
    })
    page = await browser.newPage()
    server = await startAndSeedServer()
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
