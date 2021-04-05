/* eslint-disable @typescript-eslint/no-var-requires */

const webpackBase = require('./webpack.config')

module.exports = Object.assign(webpackBase, {
  plugins: [
    ...webpackBase.plugins,
    // new CompressionPlugin({
    //   filename: '[path][base].br',
    //   algorithm: 'brotliCompress',
    //   test: /\.js$/,
    //   compressionOptions: {
    //     params: {
    //       [zlib.constants.BROTLI_PARAM_QUALITY]: 11,
    //     },
    //   },
    //   threshold: 10240,
    //   minRatio: 0.8,
    //   deleteOriginalAssets: false,
    // }),
  ],
})
