const puppeteer = require('puppeteer')
const assert = require('assert')

const catchHandler = e => {
  console.error(e)
  process.exit(1)
}

(async () => {
  const browser = await puppeteer.launch()
  const page = await browser.newPage()
  const response = await page.goto('http://localhost:6688')
  assert.ok(response.ok(), 'Should have fetched the page successfully')

  await browser.close()
})().catch(catchHandler)
