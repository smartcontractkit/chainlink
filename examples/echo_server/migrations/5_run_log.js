let chainlinkDeployer = require("../chainlink_deployer.js");
let LinkToken = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/LinkToken.sol");
let Oracle = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/Oracle.sol");
let RunLog = artifacts.require("./RunLog.sol");

module.exports = function(truffleDeployer) {
  console.log(`Create Chainlink Job`);
  chainlinkDeployer.job("only_jobid_logs_job.json", function(error, response, body) {
    console.log(`Deploying Consumer Contract with JobID ${body.id}`);
    truffleDeployer.deploy(RunLog, LinkToken.address, Oracle.address, body.id);
  }, function(error) {
    console.log("chainlink error:", error);
  });
};
