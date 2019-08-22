// truffle script

const { utils } = require('ethers')
const commandLineArgs = require('command-line-args')
const {
  optionDefinitions,
  provider,
  scriptRunner,
  wallet
} = require('./common')

const USAGE =
  'truffle exec scripts/transfer_token.js [options] <token address> <recipient address>'

const main = async () => {
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
  // send tx
  const txHash = (await wallet.sendTransaction(tx)).hash
  // wait for tx to be mined
  await provider.waitForTransaction(txHash)
  // get tx receipt
  const receipt = await provider.getTransactionReceipt(txHash)
  console.log(receipt)
  console.log(`${numTokens} transfered from ${tokenAddress} to ${recipient}`)
}

module.exports = scriptRunner(main, USAGE)
