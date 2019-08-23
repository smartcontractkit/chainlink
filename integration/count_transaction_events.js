// truffle script

const commandLineArgs = require('command-line-args')
const { abort, scriptRunner } = require('./common.js')

// compand line options
const optionDefinitions = [
  { name: 'args', type: String, multiple: true, defaultOption: true },
  { name: 'compile', type: Boolean },
  { name: 'network', type: String }
]

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
