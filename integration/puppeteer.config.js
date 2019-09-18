module.exports = {
  args: ['--no-sandbox'],
  defaultViewport: {
    width: 1366,
    height: 768,
  },
  devtools: false,
  headless: process.env.HEADLESS !== 'false',
  slowMo: process.env.SLOWMO || 10,
}
