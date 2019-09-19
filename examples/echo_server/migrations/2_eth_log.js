const EthLog = artifacts.require('EthLog')

module.exports = function(deployer) {
  deployer.deploy(EthLog)
}
