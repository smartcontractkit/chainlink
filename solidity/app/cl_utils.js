const Eth = require('ethjs')

clUtils = {
  toWei: function toWei (eth) {
    return (parseInt(eth.toString()) * 10 ** 18).toString()
  },
  getTxReceipt: function getTxReceipt (txHash) {
    return new Promise(async (resolve, reject) => {
      while (true) {
        let receipt = await clUtils.eth.getTransactionReceipt(txHash)
        if (receipt != null) {
          return resolve(receipt)
        }
      }
    })
  },
  setProvider: function setProvider (provider) {
    clUtils.provider = provider
    clUtils.eth = new Eth(provider)
  }
}

clUtils.setProvider(new Eth.HttpProvider('http://localhost:18545'))
