const puppeteer = require('puppeteer')

describe('Create Job and Run', () => {
  it('works', async () => {
    const browser = await puppeteer.launch()
    const page = await browser.newPage()
    await page.goto('http://localhost:6688')
    await expect(page).toMatch('Chainlink')
    await browser.close()
  })
})
