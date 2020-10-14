/* eslint-disable @typescript-eslint/no-var-requires */

const wp = require('@cypress/webpack-preprocessor')

module.exports = (on) => {
  const options = {
    webpackOptions: require('../../webpack.config'),
  }
  on('file:preprocessor', wp(options))

  // fix for Cypress error “The automation client disconnected. Cannot continue running tests.” when running in Docker
  // https://stackoverflow.com/a/58947968
  // ref: https://docs.cypress.io/api/plugins/browser-launch-api.html#Usage
  on('before:browser:launch', (browser = {}, args) => {
    if (browser.name === 'chrome') {
      args.push('--disable-dev-shm-usage')
      return args
    }

    return args
  })
}
