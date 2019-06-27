module.exports = {
  extends: 'react-static/.babelrc',
  presets: [
    '@babel/preset-react',
    '@babel/preset-env',
    '@babel/preset-typescript'
  ],
  plugins: [
    'react-hot-loader/babel',
    '@babel/plugin-proposal-export-namespace-from',
    '@babel/plugin-proposal-throw-expressions',
    '@babel/plugin-proposal-class-properties'
  ]
}
