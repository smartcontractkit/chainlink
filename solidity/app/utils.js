const Eth = require('ethjs')

const retries = process.env['DEPLOY_TX_CONFIRMATION_RETRIES'] || 1000
const retrySleep = process.env['DEPLOY_TX_CONFIRMATION_WAIT'] || 100

const sleep = ms => {
  return new Promise(resolve => setTimeout(resolve, ms))
}

module.exports = function Utils(provider) {
  const eth = new Eth(provider)

  return {
    eth: eth,
    provider: provider,
    toWei: eth => {
      return (parseInt(eth.toString(), 10) * 10 ** 18).toString()
    },
    getTxReceipt: txHash => {
      return new Promise(async (resolve, reject) => {
        for (let i = 0; i < retries; i++) {
          await sleep(retrySleep)

          const receipt = await eth.getTransactionReceipt(txHash)
          if (receipt != null) {
            return resolve(receipt)
          }
        }
        reject(new Error(`${txHash} unconfirmed!`))
      })
    },
    send: async params => {
      const txDefaults = { data: '' }
      return eth.sendTransaction(Object.assign(txDefaults, params))
    }
  }
}
