var Consumer = artifacts.require("./Consumer.sol");
var Oracle = artifacts.require("./Oracle.sol");
var LinkToken = artifacts.require("./LinkToken.sol");

module.exports = function(deployer) {
  deployer.deploy(Consumer, LinkToken.address, Oracle.address);
};

