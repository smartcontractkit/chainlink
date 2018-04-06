let Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(Oracle);
};
