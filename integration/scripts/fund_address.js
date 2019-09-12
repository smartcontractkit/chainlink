/* eslint-disable @typescript-eslint/no-var-requires */

// truffle script

const { utils } = require('ethers')
const commandLineArgs = require('command-line-args')
const {
  DEVNET_ADDRESS,
  optionsDefinitions,
  provider,
  scriptRunner,
  wallet,
} = require('../common.js')

const USAGE =
  'truffle exec scripts/fund_address.js [options] <optional address>'

const main = async () => {
  // parse command line args
  const options = commandLineArgs(optionsDefinitions)
  let [recipient] = options.args.slice(2)
  // transaction
  recipient = recipient || wallet.address // default
  const tx = {
    to: recipient,
    value: utils.bigNumberify(10).pow(21), // 10 ** 21
  }
  // send tx
  const devnetMinerWallet = provider.getSigner(DEVNET_ADDRESS)
  const txHash = (await devnetMinerWallet.sendTransaction(tx)).hash
  // wait for tx to be mined
  await provider.waitForTransaction(txHash)
  // get tx receipt
  const receipt = await provider.getTransactionReceipt(txHash)
  console.log(receipt)
}

module.exports = scriptRunner(main, USAGE)
