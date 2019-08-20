const { DEVNET_ADDRESS: from } = require('../common.js')
const LinkToken = artifacts.require('LinkToken')
const Oracle = artifacts.require('Oracle')

module.exports = async deployer => {
  // deploy LINK token
  await deployer.deploy(LinkToken, { from })
  const linkToken = await LinkToken.deployed()
  console.log(`Deployed LinkToken at: ${linkToken.address}`)
  // deploy Oracle
  await deployer.deploy(Oracle, linkToken.address, { from })
  const oracle = await Oracle.deployed()
  await oracle.setFulfillmentPermission(from, true, { from })
  console.log(`Deployed Oracle at: ${oracle.address}`)
}
