/* eslint-disable @typescript-eslint/no-var-requires */

const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin')
const webpack = require('webpack')

const isDev = process.env.NODE_ENV === 'development'
const extensions = ['.js', '.ts', '.tsx']

module.exports = {
  mode: isDev ? 'development' : 'production',
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: 'babel-loader',
      },
      {
        test: /\.tsx?$/,
        exclude: /node_modules/,
        use: 'babel-loader',
      },
      {
        test: /\.css$/i,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.svg$/i,
        type: 'asset/resource',
      },
    ],
  },
  resolve: {
    extensions,
    plugins: [new TsconfigPathsPlugin({ extensions })],
    fallback: {
      stream: require.resolve('stream-browserify'), // Required for @iarna/toml
    },
  },
  output: {
    filename: '[name].[contenthash].js',
    path: path.resolve(__dirname, '../core/web/assets'),
    publicPath: '/',
  },
  plugins: [
    new CleanWebpackPlugin(),
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, 'src/index.html'),
    }),
    new webpack.DefinePlugin({
      'process.env.CHAINLINK_VERSION': JSON.stringify(
        process.env.CHAINLINK_VERSION,
      ),
      'process.env.CHAINLINK_BASEURL': JSON.stringify(
        process.env.CHAINLINK_BASEURL,
      ),
    }),
  ],
}
