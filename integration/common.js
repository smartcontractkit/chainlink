const Eth = require('ethjs')
const Deployer = require('../solidity/app/deployer.js')
const Wallet = require('../solidity/app/wallet.js')
const Utils = require('../solidity/app/utils.js')

const utils = Utils(new Eth.HttpProvider('http://localhost:18545'))
const privateKey =
  '34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c'
const wallet = Wallet(privateKey, utils)

module.exports = {
  DEVNET_ADDRESS: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',

  deployer: Deployer(wallet, utils),

  abort: message => {
    return error => {
      console.error(message)
      console.error(error)
      process.exit(1)
    }
  }
}
