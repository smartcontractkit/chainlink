var DynamicConsumer = artifacts.require("./DynamicConsumer.sol");
var Chainlinked = artifacts.require("./Chainlinked.sol");
var Oracle = artifacts.require("./Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(Oracle).then(function() {
    return deployer.deploy(DynamicConsumer, Oracle.address);
  });
};

