let LINK = artifacts.require("../../../solidity/contracts/lib/LinkToken.sol");
let Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(Oracle, LINK.address);
};
