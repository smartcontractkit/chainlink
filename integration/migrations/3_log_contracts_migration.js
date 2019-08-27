const { DEVNET_ADDRESS: from } = require('../common.js')
const EthLog = artifacts.require('EthLog')
const RunLog = artifacts.require('RunLog')
const LinkToken = artifacts.require('LinkToken')
const Oracle = artifacts.require('Oracle')

module.exports = async deployer => {
  // get chainlink contracts
  const linkToken = await LinkToken.deployed()
  const oracle = await Oracle.deployed()
  // deploy EthLog contract
  await deployer.deploy(EthLog, { from })
  const ethLog = await EthLog.deployed()
  console.log(`Deployed EthLog at: ${ethLog.address}`)
  // deploy runlog contract
  await deployer.deploy(RunLog, linkToken.address, oracle.address, { from })
  const runLog = await RunLog.deployed()
  console.log(`Deployed RunLog at: ${runLog.address}`)
}
