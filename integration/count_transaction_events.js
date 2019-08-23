// truffle script

const commandLineArgs = require('command-line-args')
const { abort, scriptRunner, optionDefinitions } = require('./common.js')

const main = async () => {
  // parse command line args
  const options = commandLineArgs(optionDefinitions)
  const [txID, fromAddress] = options.args.slice(2)
  // find transaction
  const transaction = await web3.eth
    .getTransactionReceipt(txID)
    .catch(abort('Error getting transaction receipt'))
  // count events in transaction
  let count = 0
  for (let log of transaction.logs) {
    if (log.address.toLowerCase() === fromAddress.toLowerCase()) {
      count += 1
    }
  }
  console.log(`Events from ${fromAddress} in ${txID}: ${count}`)
}

module.exports = scriptRunner(main)
