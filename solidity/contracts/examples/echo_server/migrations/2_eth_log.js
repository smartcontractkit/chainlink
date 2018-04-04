var EthLog = artifacts.require("./EthLog.sol");

module.exports = function(deployer) {
  deployer.deploy(EthLog);
};
