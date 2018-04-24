let chainlinkDeployer = require("../chainlink_deployer.js");
let UptimeSLA = artifacts.require("./UptimeSLA.sol");
let Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");

module.exports = function(truffleDeployer) {
  let client = "0x542B68aE7029b7212A5223ec2867c6a94703BeE3";
  let serviceProvider = "0xB16E8460cCd76aEC437ca74891D3D358EA7d1d88";

  chainlinkDeployer.job("http_json_x10000_job.json", function(error, response, body) {
    console.log(`Deploying UptimeSLA:`)
    console.log(`\tjob: ${body.id}`);
    console.log(`\tclient: ${client}`);
    console.log(`\tservice provider: ${serviceProvider}`);

    truffleDeployer.deploy(UptimeSLA, client, serviceProvider, Oracle.address, body.id, {
      value: 1000000000
    });
  }, function(error) {
    console.log("chainlink error:", error);
  });
};
