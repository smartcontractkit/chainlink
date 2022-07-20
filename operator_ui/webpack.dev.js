/* eslint-disable @typescript-eslint/no-var-requires */
const webpack = require('webpack')
const webpackBase = require('./webpack.config')

module.exports = Object.assign(webpackBase, {
  devtool: 'inline-source-map',
  devServer: {
    port: 3000,
    contentBase: './dist',
    historyApiFallback: true,
    hot: true,
  },
  plugins: [...webpackBase.plugins, new webpack.HotModuleReplacementPlugin()],
})
