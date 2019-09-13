import pupExpect from 'expect-puppeteer'

const clickNavigationTag = tagName => async (page, text) => {
  // XXX: Some buttons/links don't do anything if you click them too quickly,
  // so for, now, add a small delay
  await page.waitFor(500)
  await pupExpect(page).toClick(tagName, { text })
}

module.exports = {
  clickLink: clickNavigationTag('a'),
  clickButton: clickNavigationTag('button'),

  consoleLogger: page => msg => {
    console.log(`PAGE LOG url: ${page.url()} | msg: ${msg.text()}`)
  },

  clickTransactionsMenuItem: async page => {
    return pupExpect(page).toClick('li > a', { text: 'Transactions' })
  },

  signIn: async (page, email, password) => {
    await pupExpect(page).toFill('form input[id=email]', email)
    await pupExpect(page).toFill('form input[id=password]', 'twochains')
    return pupExpect(page).toClick('form button')
  },
}
