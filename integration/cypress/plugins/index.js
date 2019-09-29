/* eslint-disable @typescript-eslint/no-var-requires */

const wp = require('@cypress/webpack-preprocessor')

module.exports = on => {
  const options = {
    webpackOptions: require('../../webpack.config'),
  }
  on('file:preprocessor', wp(options))
}
