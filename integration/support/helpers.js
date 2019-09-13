const clickLink = async (page, title) => {
  // XXX: This buttons don't do anything if you click them too quickly, so for
  // now, add a small delay
  await page.waitFor(500)
  await expect(page).toClick('a', { text: title })
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
    return expect(page).toClick('li > a', { text: 'Transactions' })
  },

  signIn: async (page, email) => {
    await expect(page).toFill('form input[id=email]', email)
    await expect(page).toFill('form input[id=password]', 'twochains')
    return expect(page).toClick('form button')
  },
}
