const HtmlWebpackPlugin = require('html-webpack-plugin')
const DynamicCdnWebpackPlugin = require('dynamic-cdn-webpack-plugin')
const CompressionPlugin = require('compression-webpack-plugin')
const BrotliPlugin = require('brotli-webpack-plugin')
const webpack = require('webpack')
const clientPkg = require('./package.json')
const serverPkg = require('../package.json')
const GitRevisionPlugin = require('git-revision-webpack-plugin')
const gitRevisionPlugin = new GitRevisionPlugin({ branch: true })

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
      new webpack.DefinePlugin({
        __EXPLORER_CLIENT_VERSION__: JSON.stringify(clientPkg.version),
        __EXPLORER_SERVER_VERSION__: JSON.stringify(serverPkg.version),
        __GIT_SHA__: JSON.stringify(gitRevisionPlugin.commithash()),
        __GIT_BRANCH__: JSON.stringify(gitRevisionPlugin.branch()),
        __REACT_APP_EXPLORER_GA_ID__: JSON.stringify(process.env.REACT_APP_EXPLORER_GA_ID)
      })
    ],
  },
  eslint: {
    enable: false
  },
}
