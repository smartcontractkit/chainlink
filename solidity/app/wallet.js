const Tx = require('ethereumjs-tx')
const EthWallet = require('ethereumjs-wallet')

module.exports = function Wallet(key, utils) {
  const privateKey = Buffer.from(key, 'hex')
  const wallet = EthWallet.fromPrivateKey(privateKey)
  const address = `0x${wallet.getAddress().toString('hex')}`
  const eth = utils.eth
  const nextNonce = async () => eth.getTransactionCount(address, 'pending')

  return {
    address: address,
    nextNonce: nextNonce,
    send: async params => {
      const defaults = {
        nonce: await nextNonce(),
        chainId: 0
      }
      let tx = new Tx(Object.assign(defaults, params))
      tx.sign(privateKey)
      let txHex = tx.serialize().toString('hex')
      return eth.sendRawTransaction(txHex)
    },
    call: async params => {
      return eth.call(params)
    }
  }
}
