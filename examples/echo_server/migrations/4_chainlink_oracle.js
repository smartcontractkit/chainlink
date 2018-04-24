let LinkToken = artifacts.require("../../../solidity/contracts/LinkToken.sol");
let Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(Oracle, LinkToken.address);
};
