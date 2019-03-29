const Deployer = require('./deployer.js')
const Eth = require('ethjs')
const Utils = require('./utils.js')
const Wallet = require('./wallet.js')
const Web3 = require('web3')

const privateKey =
  process.env['PRIVATE_KEY'] ||
  '4d6cf3ce1ac71e79aa33cf481dedf2e73acb548b1294a70447c960784302d2fb'
const providerURL = process.env['ETH_HTTP_URL'] || 'http://localhost:18545'
const utils = new Utils(new Eth.HttpProvider(providerURL))
const wallet = new Wallet(privateKey, utils)
const deployer = Deployer(wallet, utils)
const web3 = new Web3(providerURL)

module.exports = {
  abi: require('ethereumjs-abi'),
  deployer: deployer,
  utils: utils,
  wallet: wallet,
  web3: web3
}
