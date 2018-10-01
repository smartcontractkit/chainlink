let LINK = artifacts.require("link_token/contracts/LinkToken.sol");
let Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(Oracle, LINK.address);
};
