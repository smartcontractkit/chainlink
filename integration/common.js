const Eth = require('ethjs')
const Deployer = require('../evm/app/deployer.js')
const Wallet = require('../evm/app/wallet.js')
const Utils = require('../evm/app/utils.js')

const port = process.env.ETH_HTTP_PORT || `18545`
const utils = Utils(new Eth.HttpProvider(`http://localhost:${port}`))
const privateKey =
  '34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c'
const wallet = Wallet(privateKey, utils)

module.exports = {
  DEVNET_ADDRESS: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',

  deployer: Deployer(wallet, utils),
  port: port,

  abort: message => {
    return error => {
      console.error(message)
      console.error(error)
      process.exit(1)
    }
  }
}
