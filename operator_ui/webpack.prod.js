/* eslint-disable @typescript-eslint/no-var-requires */
const CompressionPlugin = require('compression-webpack-plugin')
const webpackBase = require('./webpack.config')

module.exports = Object.assign(webpackBase, {
  output: {
    ...webpackBase.output,
    publicPath: '/assets/', // JS files are served from `/assets` by web
  },
  plugins: [...webpackBase.plugins, new CompressionPlugin({})],
})
