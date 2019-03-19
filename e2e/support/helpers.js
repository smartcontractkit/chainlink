module.exports = {
  clickNewJobButton: async page => {
    await page.waitFor(500)
    return expect(page).toClick('a', { text: 'New Job' })
  },

  clickTransactionsMenuItem: async page => {
    return expect(page).toClick('li > a', { text: 'Transactions' })
  },

  signIn: async (page, email, password) => {
    await expect(page).toFill('form input[id=email]', email)
    await expect(page).toFill('form input[id=password]', 'twochains')
    return expect(page).toClick('form button')
  }
}
