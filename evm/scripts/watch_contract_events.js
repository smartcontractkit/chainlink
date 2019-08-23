// truffle script

const ethers = require('ethers')
const commandLineArgs = require('command-line-args')
const { optionsDefinitions, wallet } = require('./common')

const USAGE =
  'truffle exec scripts/watch_contract_events.js [options] <contract name> <contract address>'

const main = async () => {
  // parse command line args
  const options = commandLineArgs(optionsDefinitions)
  const [contractName, contractAddress] = options.args.slice(2)
  // import abi & bytecode from build
  const { abi, bytecode } = artifacts.require(contractName)
  // watch events
  console.log(`Watching events at ${contractAddress}`)
  const contractFactory = new ethers.ContractFactory(abi, bytecode, wallet)
  const contract = contractFactory.attach(contractAddress)
  contract.on('*', console.log)
}

module.exports = scriptRunner(main, USAGE)
