const pupExpect = require('expect-puppeteer')

const clickLink = async (page, title) => {
  // XXX: Some buttons don't do anything if you click them too quickly,
  // so for, now, add a small delay
  await page.waitFor(500)
  await pupExpect(page).toClick('a', { text: title })
  await page.waitForNavigation({
    waitUntil: 'networkidle0',
  })
}

module.exports = {
  clickLink: clickLink,

  consoleLogger: page => msg => {
    console.log(`PAGE LOG url: ${page.url()} | msg: ${msg.text()}`)
  },

  clickNewJobButton: async page => clickLink(page, 'New Job'),

  clickBridgesTab: async page => clickLink(page, 'Bridges'),

  clickNewBridgeButton: async page => clickLink(page, 'New Bridge'),

  clickTransactionsMenuItem: async page => {
    return pupExpect(page).toClick('li > a', { text: 'Transactions' })
  },

  signIn: async (page, email, password) => {
    await pupExpect(page).toFill('form input[id=email]', email)
    await pupExpect(page).toFill('form input[id=password]', 'twochains')
    return pupExpect(page).toClick('form button')
  },
}
