const Tx = require('ethereumjs-tx')
const EthWallet = require('ethereumjs-wallet')

module.exports = function Wallet (key, utils) {
  const privateKey = Buffer.from(key, 'hex')
  const wallet = EthWallet.fromPrivateKey(privateKey)
  const address = wallet.getAddress().toString('hex')
  const eth = utils.eth

  this.address = address
  this.send = async (params) => {
    const defaults = {
      nonce: await this.nextNonce(),
      chainId: 0
    }
    let tx = new Tx(Object.assign(defaults, params))
    tx.sign(privateKey)
    let txHex = tx.serialize().toString('hex')
    return eth.sendRawTransaction(txHex)
  }
  this.nextNonce = () => {
    return eth.getTransactionCount(address)
  }
}
