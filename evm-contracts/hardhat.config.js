require('@nomiclabs/hardhat-truffle5')
require('@nomiclabs/hardhat-web3')
require('hardhat-gas-reporter')

// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task('accounts', 'Prints the list of accounts', async (args, hre) => {
  console.log(await hre.web3.eth.getAccounts())
})

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: {
    version: '0.8.4',
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  paths: {
    sources: './src/v0.8',
    tests: './test/v0.8',
    cache: './cache/v0.8',
    artifacts: './artifacts/v0.8',
  },
  mocha: {
    timeout: 20000,
  },
}
