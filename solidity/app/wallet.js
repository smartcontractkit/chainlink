const Tx = require('ethereumjs-tx')
const EthWallet = require('ethereumjs-wallet')

module.exports = function Wallet (key, utils) {
  this.privateKey = Buffer.from(key, 'hex')
  const wallet = EthWallet.fromPrivateKey(this.privateKey)
  this.address = wallet.getAddress().toString('hex')

  this.send = async function (params) {
    let eth = utils.eth
    let defaults = {
      nonce: await eth.getTransactionCount(this.address),
      chainId: 0
    }
    let tx = new Tx(Object.assign(defaults, params))
    tx.sign(this.privateKey)
    let txHex = tx.serialize().toString('hex')
    return eth.sendRawTransaction(txHex)
  }
}
