// truffle script

const ethers = require('ethers')
const commandLineArgs = require('command-line-args')
const { wallet } = require('../chainlink.config')

// compand line options
const optionDefinitions = [
  { name: 'args', type: String, multiple: true, defaultOption: true },
  { name: 'compile', type: Boolean },
  { name: 'network', type: String }
]

module.exports = async function(callback) {
  // parse command line args
  const options = commandLineArgs(optionDefinitions)
  const [contractName, contractAddress] = options.args.slice(2)
  try {
    // import abi & bytecode from build
    const { abi, bytecode } = artifacts.require(contractName)
    // watch events
    console.log(`Watching events at ${contractAddress}`)
    const contractFactory = new ethers.ContractFactory(abi, bytecode, wallet)
    const contract = contractFactory.attach(contractAddress)
    contract.on('*', console.log)
  } catch (error) {
    console.error('Usage: truffle exec scripts/watch_contract_events.js [options] ' +
    '<contract name> <contract address>')
    callback(error)
  }
}
