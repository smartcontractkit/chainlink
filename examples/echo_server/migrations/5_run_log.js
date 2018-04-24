let LinkToken = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/LinkToken.sol");
let Oracle = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/Oracle.sol");
let RunLog = artifacts.require("./RunLog.sol");

module.exports = function(deployer) {
  deployer.deploy(RunLog, LinkToken.address, Oracle.address);
  LinkToken.deployed().then(function(instance) {
    return instance.transfer(RunLog.address, 5000000000000000000);
  }).then(console.log, console.log);
};
