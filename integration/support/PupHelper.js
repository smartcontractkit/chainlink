/* eslint-disable @typescript-eslint/no-var-requires */

const pupExpect = require('expect-puppeteer')

module.exports = class PupHelper {
  constructor(page) {
    this.page = page
    // page.setDefaultTimeout(29000)
    this.page.on('console', msg => {
      console.log(`PAGE LOG url: ${page.url()} | msg: ${msg.text()}`)
    })
  }

  async clickButton(text) {
    await this.page.waitFor(500)
    await this.waitForContent('button', text)
    await pupExpect(this.page).toClick('button', { text })
  }

  async clickLink(text) {
    await this.page.waitFor(500)
    await this.waitForContent('a', text)
    await pupExpect(this.page).toClick('a', { text })
  }

  async signIn(email = 'notreal@fakeemail.ch', password = 'twochains') {
    await pupExpect(this.page).toFill('form input[id=email]', email)
    await pupExpect(this.page).toFill('form input[id=password]', password)
    return pupExpect(this.page).toClick('form button')
  }

  async waitForContent(tagName, content) {
    const xpath = `//${tagName}[contains(., '${content}')]`
    try {
      return await this.page.waitForXPath(xpath)
    } catch {
      throw `Unable to find <${tagName}> tag with content: '${content}'`
    }
  }

  async waitForNotification(notification) {
    return await this.waitForContent('p', notification)
  }
}
