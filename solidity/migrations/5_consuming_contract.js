let chainlinkDeployer = require("../chainlink_deployer.js");
let Consumer = artifacts.require("./Consumer.sol");
let Oracle = artifacts.require("./Oracle.sol");
let LinkToken = artifacts.require("./LinkToken.sol");
let fs = require('fs');
let jobJson = '{"initiators":[{"type":"runlog"}],"tasks":[{"type":"httpGet"},{"type":"jsonParse"},{"type":"multiply","times":100},{"type":"ethuint256"},{"type":"ethtx"}]}';

module.exports = function(truffleDeployer) {
  chainlinkDeployer.job(jobJson, function(error, response, body) {
    console.log(`Deploying Consumer:`)
    console.log(`\tjob: ${body.id}`);
    truffleDeployer.deploy(Consumer, LinkToken.address, Oracle.address, body.id);
  }, function(error) {
    console.log("chainlink error:", error);
  });
};
