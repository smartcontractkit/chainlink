const clWallet = require('./cl_wallet.js')

const Eth = require('ethjs')

global.clUtils = global.clUtils || {
  personalAccount: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
  toWei: function toWei (eth) {
    return (parseInt(eth.toString(), 10) * 10 ** 18).toString()
  },
  getTxReceipt: function getTxReceipt (txHash) {
    return new Promise(async (resolve, reject) => {
      for (let i = 0; i < 1000; i++) {
        let receipt = await clUtils.eth.getTransactionReceipt(txHash)
        if (receipt != null) {
          return resolve(receipt)
        }
      }
      reject(`${txHash} unconfirmed!`)
    })
  },
  setProvider: function setProvider (provider) {
    clUtils.provider = provider
    clUtils.eth = new Eth(provider)
  },
  send: async function send (params) {
    let defaults = {
      data: '',
      from: clUtils.personalAccount
    }
    return clUtils.eth.sendTransaction(Object.assign(defaults, params))
  }
}
clUtils.setProvider(new Eth.HttpProvider('http://localhost:18545'))

module.exports = global.clUtils
