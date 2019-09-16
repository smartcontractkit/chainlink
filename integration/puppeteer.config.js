module.exports = {
  devtools: false,
  headless: process.env.HEADLESS !== 'false',
  args: ['--no-sandbox'],
  defaultViewport: {
    width: 1366,
    height: 768,
  },
}
