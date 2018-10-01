let LINK = artifacts.require("link_token/contracts/LinkToken.sol");

module.exports = function(deployer) {
  deployer.deploy(LINK);
};

