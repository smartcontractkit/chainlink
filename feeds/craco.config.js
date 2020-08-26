const HtmlWebpackPlugin = require('html-webpack-plugin')
const DynamicCdnWebpackPlugin = require('dynamic-cdn-webpack-plugin')
const CompressionPlugin = require('compression-webpack-plugin')
const BrotliPlugin = require('brotli-webpack-plugin')
const webpack = require('webpack')

module.exports = {
  webpack: {
    plugins: [
      new HtmlWebpackPlugin(),
      new DynamicCdnWebpackPlugin(),
      new CompressionPlugin({
        test: /\.js$|\.css$|\.html$|\.svg/,
        filename: '[path].gz[query]',
        algorithm: 'gzip',
        threshold: 0,
        minRatio: 0.8,
      }),
      new BrotliPlugin({
        test: /\.(js|css|html|svg)$/,
        asset: '[path].br[query]',
        threshold: 0,
        minRatio: 0.8,
      }),
    ],
  },
  eslint: {
    enable: false
  },
}
