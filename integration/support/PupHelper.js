const pupExpect = require('expect-puppeteer')

module.exports = class PupHelper {
  constructor(page) {
    this.page = page
    // page.setDefaultTimeout(3000)
    this.page.on('console', msg => {
      console.log(`PAGE LOG url: ${page.url()} | msg: ${msg.text()}`)
    })
  }

  async clickButton(text) {
    // XXX: Some buttons/links don't do anything if you click them too quickly,
    // so for, now, add a small delay
    await this.page.waitFor(500)
    await pupExpect(this.page).toClick('button', { text })
  }

  async clickLink(text, selector = 'a') {
    return Promise.all([
      this.page.waitForNavigation(),
      pupExpect(this.page).toClick(selector, { text }),
    ])
  }

  async clickMenuItem(text) {
    await this.clickLink(text, 'li > a')
  }

  async signIn(email = 'notreal@fakeemail.ch', password = 'twochains') {
    await pupExpect(this.page).toFill('form input[id=email]', email)
    await pupExpect(this.page).toFill('form input[id=password]', password)
    return pupExpect(this.page).toClick('form button')
  }

  async waitForContent(tagName, content) {
    const xpath = `//${tagName}[contains(text(), '${content}')]`
    try {
      return await this.page.waitForXPath(xpath)
    } catch {
      throw `Unable to find <${tagName}> tag with content: "${content}"`
    }
  }

  async waitForNotification(notification) {
    return await this.waitForContent('p', notification)
  }
}
