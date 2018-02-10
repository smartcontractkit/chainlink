var DynamicConsumer = artifacts.require("./DynamicConsumer.sol");
var Oracle = artifacts.require("./Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(DynamicConsumer, Oracle.address);
};

