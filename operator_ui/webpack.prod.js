/* eslint-disable @typescript-eslint/no-var-requires */
const CompressionPlugin = require('compression-webpack-plugin')
const webpackBase = require('./webpack.config')

module.exports = Object.assign(webpackBase, {
  plugins: [...webpackBase.plugins, new CompressionPlugin({})],
})
