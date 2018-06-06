const clUtils = require('./cl_utils.js')

const Tx = require('ethereumjs-tx')
const Wallet = require('ethereumjs-wallet')

global.clWallet = global.clWallet || {
  send: async function (params) {
    let eth = clUtils.eth
    let defaults = {
      nonce: await eth.getTransactionCount(this.address),
      chainId: 0
    }
    let tx = new Tx(Object.assign(defaults, params))
    tx.sign(this.privateKey)
    let txHex = tx.serialize().toString('hex')
    return eth.sendRawTransaction(txHex)
  },
  setDefaultKey: function (key) {
    this.privateKey = Buffer.from(key, 'hex')
    const wallet = Wallet.fromPrivateKey(this.privateKey)
    this.address = wallet.getAddress().toString('hex')
  }
}
clWallet.setDefaultKey('4d6cf3ce1ac71e79aa33cf481dedf2e73acb548b1294a70447c960784302d2fb')

module.exports = global.clWallet
