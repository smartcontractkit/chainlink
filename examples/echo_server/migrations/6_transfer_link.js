const clmigration = require('../clmigration.js')
const LinkToken = artifacts.require('LinkToken')
const RunLog = artifacts.require('RunLog')
const devnetAddress = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'

module.exports = clmigration(function(deployer) {
  LinkToken.deployed().then(async function(linkInstance) {
    await RunLog.deployed()
      .then(async function(runLogInstance) {
        await linkInstance.transfer(
          runLogInstance.address,
          web3.utils.toWei('1000')
        )
        await linkInstance.transfer(devnetAddress, web3.utils.toWei('1000'))
      })
      .catch(console.log)
  })
})
