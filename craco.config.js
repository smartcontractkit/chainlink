const HtmlWebpackPlugin = require('html-webpack-plugin')
const DynamicCdnWebpackPlugin = require('dynamic-cdn-webpack-plugin')

module.exports = {
  webpack: {
    plugins: [new HtmlWebpackPlugin(), new DynamicCdnWebpackPlugin()]
  }
}
