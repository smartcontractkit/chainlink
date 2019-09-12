const GetMoney = artifacts.require('./GetMoney.sol')

module.exports = function(deployer) {
  deployer.deploy(GetMoney)
}
