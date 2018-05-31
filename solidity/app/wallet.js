const Tx = require('ethereumjs-tx')
const Wallet = require('ethereumjs-wallet')

let privateKey = Buffer.from('4d6cf3ce1ac71e79aa33cf481dedf2e73acb548b1294a70447c960784302d2fb', 'hex')
let wallet = Wallet.fromPrivateKey(privateKey)
let address = wallet.getAddress().toString('hex')

module.exports = {
  address: address,
  privateKey: privateKey,
  send: async function send (params) {
    let eth = clUtils.eth
    let defaults = {
      nonce: await eth.getTransactionCount(address),
      chainId: 0
    }
    let tx = new Tx(Object.assign(defaults, params))
    tx.sign(privateKey)
    let txHex = tx.serialize().toString('hex')
    return eth.sendRawTransaction(txHex)
  }
}
