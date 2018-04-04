let Oracle = artifacts.require("../../../Oracle.sol");
let RunLog = artifacts.require("./RunLog.sol");

module.exports = function(deployer) {
  deployer.deploy(RunLog, Oracle.address);
};
