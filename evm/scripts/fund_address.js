// truffle script

const { utils } = require('ethers')
const commandLineArgs = require('command-line-args')
const { wallet, provider, devnetMiner } = require('../chainlink.config')

// compand line options
const optionDefinitions = [
  { name: 'args', type: String, multiple: true, defaultOption: true },
  { name: 'compile', type: Boolean },
  { name: 'network', type: String }
]

module.exports = async function(callback) {
  // parse command line args
  const options = commandLineArgs(optionDefinitions)
  let [recipient] = options.args.slice(2)
  // transaction
  recipient = recipient || wallet.address // default
  const tx = {
    to: recipient,
    value: utils.bigNumberify(10).pow(21) // 10 ** 21
  }
  try {
    // send tx
    const devnetMinerWallet = provider.getSigner(devnetMiner)
    const txHash = (await devnetMinerWallet.sendTransaction(tx)).hash
    // wait for tx to be mined
    await provider.waitForTransaction(txHash)
    // get tx receipt
    const receipt = await provider.getTransactionReceipt(txHash)
    console.log(receipt)
    callback()
  } catch (error) {
    console.error('Usage: truffle exec scripts/fund_dev_wallet.js [options] ' +
    '<optional address>')
    callback(error)
  }
}
