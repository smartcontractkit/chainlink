// truffle script

const ethers = require('ethers')
const commandLineArgs = require('command-line-args')
const { optionsDefinitions, scriptRunner, wallet } = require('../common')

const USAGE =
  'truffle exec scripts/deploy.js [options] <contract name> <constructor args...>'

const main = async () => {
  // parse command line args
  const options = commandLineArgs(optionsDefinitions)
  const [contractName, ...constructorArgs] = options.args.slice(2)
  // import abi & bytecode from build
  const { abi, bytecode } = artifacts.require(contractName)
  // deploy
  const contractFactory = new ethers.ContractFactory(abi, bytecode, wallet)
  const contract = await contractFactory.deploy(...constructorArgs)
  console.log(`${contractName} contract successfully deployed at: ${contract.address}`)
}

module.exports = scriptRunner(main, USAGE)
