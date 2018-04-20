var Consumer = artifacts.require("./Consumer.sol");
var Oracle = artifacts.require("./Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(Consumer, Oracle.address);
};

