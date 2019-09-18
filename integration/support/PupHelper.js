/* eslint-disable @typescript-eslint/no-var-requires */

const pupExpect = require('expect-puppeteer')

module.exports = class PupHelper {
  constructor(page) {
    this.page = page
    page.setDefaultTimeout(3000)
    this.page.on('console', msg => {
      console.log(`PAGE LOG url: ${page.url()} | msg: ${msg.text()}`)
    })
  }

  async clickButton(content) {
    await this.waitForContent('button', content)
    await this.nativeClick('button', content)
  }

  async clickLink(content) {
    await this.waitForContent('a', content)
    await this.nativeClick('a', content)
  }

  // using puppeteer's #click method doesn't reliably trigger navigation
  // workaround is to trigger click natively
  async nativeClick(...params) {
    await this.page.evaluate((tagName, content) => {
      const tags = Array.from(document.querySelectorAll(tagName))
      tags.find(tag => tag.innerText.includes(content)).click()
    }, ...params)
  }

  async signIn(email = 'notreal@fakeemail.ch', password = 'twochains') {
    await pupExpect(this.page).toFill('form input[id=email]', email)
    await pupExpect(this.page).toFill('form input[id=password]', password)
    await Promise.all([
      pupExpect(this.page).toClick('form button'),
      this.page.waitForNavigation(),
    ])
  }

  async waitForContent(tagName, content) {
    const xpath = this._xpath(tagName, content)
    try {
      return await this.page.waitForXPath(xpath)
    } catch {
      throw `Unable to find <${tagName}> tag with content: '${content}'`
    }
  }

  // async refreshUntilContentPresent(tagName, content, timeout = 5000) {
  //   const xpath = this._xpath(tagName, content)
  //   return Promise.race([
  //     setTimeout(() => {
  //       throw `Unable to find <${tagName}> tag with content: '${content}'`
  //     }, timeout),
  //     async () => {
  //       let contentPresent = false
  //       while (!contentPresent) {
  //         await this.page.reload()
  //         contentPresent = (await this.page.$x(xpath)).length != 0
  //       }
  //     },
  //   ])
  // }

  async waitForNotification(notification) {
    return await this.waitForContent('p', notification)
  }

  _xpath(tagName, content) {
    // searches for tag with content any # of children deep
    return `//${tagName}[contains(., '${content}')]`
  }
}
