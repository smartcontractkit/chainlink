const clickLink = async (page, title) => {
  return expect(page).toClick('a', { text: title })
}

module.exports = {
  clickLink: clickLink,

  clickNewJobButton: async page => {
    // XXX: This button doesn't do anything if you click it too quickly, so for
    // now, add a small delay
    await page.waitFor(500)
    return clickLink(page, 'New Job')
  },

  clickTransactionsMenuItem: async page => {
    return expect(page).toClick('li > a', { text: 'Transactions' })
  },

  signIn: async (page, email, password) => {
    await expect(page).toFill('form input[id=email]', email)
    await expect(page).toFill('form input[id=password]', 'twochains')
    return expect(page).toClick('form button')
  },
}
