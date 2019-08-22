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
  const [contractName, ...constructorArgs] = options.args.slice(2)
  try {
    // import abi & bytecode from build
    const { abi, bytecode } = artifacts.require(contractName)
    // deploy
    const contractFactory = new ethers.ContractFactory(abi, bytecode, wallet)
    const contract = await contractFactory.deploy(...constructorArgs)
    console.log(`${contractName} contract successfully deployed at: ${contract.address}`)
    callback()
  } catch (error) {
    console.error('Usage: truffle exec scripts/deploy.js [options] ' +
    '<contract name> <constructor args...>')
    callback(error)
  }
}
