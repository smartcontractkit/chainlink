const waitWithTimeout = async (promise, taskName, timeout) => {
  let rejectCallback
  const timeoutError = new Error(
    `waiting for ${taskName} failed: timeout ${timeout}ms exceeded`
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
    checkPage = async () => {
      const content = await page.content()
      let match = content
        .replace(/\s+/g, ' ')
        .trim()
        .match(regex)
      if (match) {
        return match
      }
      await page.reload()
      return await checkPage(page, regex)
    }
    let match = await waitWithTimeout(checkPage, 'scrape', 30000)
    return await match()
  }
}
