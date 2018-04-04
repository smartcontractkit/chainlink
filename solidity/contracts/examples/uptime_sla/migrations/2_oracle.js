let Oracle = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/Oracle.sol");

module.exports = function(deployer) {
  deployer.deploy(Oracle);
};
