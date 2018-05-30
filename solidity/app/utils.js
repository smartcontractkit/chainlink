const Eth = require('ethjs')

let eth = new Eth(new Eth.HttpProvider('http://localhost:18545'))

module.exports = {
  toWei: function toWei (eth) {
    return (parseInt(eth.toString()) * 10 ** 18).toString()
  },
  getTxReceipt: function getTxReceipt (txHash) {
    return new Promise(async (resolve, reject) => {
      while (true) {
        let receipt = await eth.getTransactionReceipt(txHash)
        if (receipt != null) {
          return resolve(receipt)
        }
      }
    })
  }
}
