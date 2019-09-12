module.exports = {
  root: true,
  extends: [
    // TODO: might need to clean this up since truffle uses mocha and we're also using jest here
    '@chainlink/eslint-config/node',
    '@chainlink/eslint-config/jest',
    '@chainlink/eslint-config/truffle',
  ],
}
