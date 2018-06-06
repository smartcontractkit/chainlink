const Eth = require('ethjs')

const retries = process.env['DEPLOY_TX_CONFIRMATION_RETRIES'] || 1000
const retrySleep = process.env['DEPLOY_TX_CONFIRMATION_WAIT'] || 100

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

global.clUtils = global.clUtils || {
  personalAccount: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
  toWei: function (eth) {
    return (parseInt(eth.toString(), 10) * 10 ** 18).toString()
  },
  getTxReceipt: function (txHash) {
    return new Promise(async (resolve, reject) => {
      for (let i = 0; i < retries; i++) {
        await sleep(retrySleep)

        const receipt = await clUtils.eth.getTransactionReceipt(txHash)
        if (receipt != null) {
          return resolve(receipt)
        }
      }
      reject(`${txHash} unconfirmed!`)
    })
  },
  setProvider: function (provider) {
    this.provider = provider
    this.eth = new Eth(provider)
  },
  send: async function (params) {
    const defaults = {
      data: '',
      from: this.personalAccount
    }
    return this.eth.sendTransaction(Object.assign(defaults, params))
  }
}
clUtils.setProvider(new Eth.HttpProvider('http://localhost:18545'))

module.exports = global.clUtils
