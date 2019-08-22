// truffle script

const { utils } = require('ethers')
const commandLineArgs = require('command-line-args')
const { wallet, provider } = require('../chainlink.config')

// compand line options
const optionDefinitions = [
  { name: 'args', type: String, multiple: true, defaultOption: true },
  { name: 'compile', type: Boolean },
  { name: 'network', type: String }
]

module.exports = async function(callback) {
  // parse command line args
  const options = commandLineArgs(optionDefinitions)
  let [tokenAddress, recipient] = options.args.slice(2)
  // encode function call
  const funcSelector = '0xa9059cbb' // "transfer(address,uint256)"
  const numTokens = utils.bigNumberify(10).pow(21)
  const encodedParams = utils.defaultAbiCoder.encode(['address', 'uint256'], [recipient, numTokens])
  const data = utils.hexlify(utils.concat([funcSelector, encodedParams]))
  // transaction
  const tx = {
    data,
    to: tokenAddress
  }
  try {
    // send tx
    const txHash = (await wallet.sendTransaction(tx)).hash
    // wait for tx to be mined
    await provider.waitForTransaction(txHash)
    // get tx receipt
    const receipt = await provider.getTransactionReceipt(txHash)
    console.log(receipt)
    console.log(`${numTokens} transfered from ${tokenAddress} to ${recipient}`)
    callback()
  } catch (error) {
    console.error('Usage: truffle exec scripts/transfer_token.js [options] ' +
    '<token address> <recipient address>')
    callback(error)
  }
}
