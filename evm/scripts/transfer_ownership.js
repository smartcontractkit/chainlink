// truffle script

const { utils } = require('ethers')
const commandLineArgs = require('command-line-args')
const {
  optionsDefinitions,
  provider,
  scriptRunner,
  wallet
} = require('../common')

const USAGE =
  'truffle exec scripts/transfer_owner.js [options] <owned address> <recipient address>'

const main = async () => {
  // parse command line args
  const options = commandLineArgs(optionsDefinitions)
  let [owned, recipient] = options.args.slice(2)
  // encode function call
  const funcSelector = '0xf2fde38b' // "transferOwnership(address)"
  const encodedParams = utils.defaultAbiCoder.encode(['address'], [recipient])
  const data = utils.hexlify(utils.concat([funcSelector, encodedParams]))
  // transaction
  const tx = { data, to: owned }
  // send tx
  const txHash = (await wallet.sendTransaction(tx)).hash
  // wait for tx to be mined
  await provider.waitForTransaction(txHash)
  // get tx receipt
  const receipt = await provider.getTransactionReceipt(txHash)
  console.log(receipt)
  console.log(`ownership of ${owned} transferred to ${recipient}`)
}

module.exports = scriptRunner(main, USAGE)
