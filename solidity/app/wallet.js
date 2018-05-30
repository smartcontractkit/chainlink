const Eth = require('ethjs')
const Tx = require('ethereumjs-tx')
const Wallet = require('ethereumjs-wallet')

let eth = clUtils.eth
let privateKey = Buffer.from('4d6cf3ce1ac71e79aa33cf481dedf2e73acb548b1294a70447c960784302d2fb', 'hex')
let wallet = Wallet.fromPrivateKey(privateKey)
let address = wallet.getAddress().toString('hex')


module.exports = {
  address: address,
  privateKey: privateKey,
  send: async function send (params) {
    let defaults = {nonce: await eth.getTransactionCount(address)}
    let tx = new Tx(Object.assign(defaults, params))
    tx.sign(privateKey)
    return eth.sendRawTransaction(tx.serialize().toString('hex'))
  }
}
