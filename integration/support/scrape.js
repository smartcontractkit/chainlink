const waitWithTimeout = async (promise, taskName, timeout) => {
  let rejectCallback
  const timeoutError = new Error(
    `waiting for ${taskName} failed: timeout ${timeout}ms exceeded`,
  )
  const timeoutPromise = new Promise((resolve, reject) => {
    rejectCallback = reject
  })
  const timeoutTimer = setTimeout(() => rejectCallback(timeoutError), timeout)
  try {
    return await Promise.race([promise, timeoutPromise])
  } finally {
    clearTimeout(timeoutTimer)
  }
}

module.exports = {
  // Scrape matches against the page content and returns the match groups found
  // if any, it will refresh the page until the match is found or timeout
  // occurs.
  scrape: async (page, regex) => {
    const checkPage = async () => {
      const element = await page.$('body')
      const text = await (await element.getProperty('textContent')).jsonValue()
      const match = text
        .replace(/\s+/g, ' ')
        .trim()
        .match(regex)
      if (match) {
        return match
      }
      await page.reload()
      return await checkPage(page, regex)
    }
    const match = await waitWithTimeout(checkPage, 'scrape', 30000)
    expect(match).toBeDefined()
    return (await match())[1]
  },
}
