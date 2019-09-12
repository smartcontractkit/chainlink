module.exports = {
  extends: ['@chainlink/eslint-config/mocha'],
  globals: {
    assert: 'readonly',
    artifacts: 'readonly',
    web3: 'readonly',
    contract: 'readonly',
  },
}
