let LinkToken = artifacts.require("../../../solidity/contracts/LinkToken.sol");

module.exports = function(deployer) {
  deployer.deploy(LinkToken);
};
