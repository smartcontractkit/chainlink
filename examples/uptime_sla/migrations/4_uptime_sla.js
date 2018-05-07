let chainlinkDeployer = require("../chainlink_deployer.js");
let UptimeSLA = artifacts.require("./UptimeSLA.sol");
let Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");
let LINK = artifacts.require("../../../solidity/contracts/LinkToken.sol");

let url = "http://chainlink:twochains@localhost:6688/v2/specs";
let job = {
  "_comment": "GETs a number from JSON, multiplies by 10,000, and reports uint256",
  "initiators": [
    { "type": "runlog"}
  ],
  "tasks": [
    {"type": "httpGet"},
    {"type": "jsonParse"},
    {"type": "multiply", "times": 10000},
    {"type": "ethuint256"},
    {"type": "ethtx"}
  ]
}

module.exports = function(truffleDeployer) {
  truffleDeployer.then(async () => {
    let client = "0x542B68aE7029b7212A5223ec2867c6a94703BeE3";
    let serviceProvider = "0xB16E8460cCd76aEC437ca74891D3D358EA7d1d88";

    let body = await chainlinkDeployer.job(url, job);
    console.log(`Deploying UptimeSLA:`)
    console.log(`\tjob: ${body.id}`);
    console.log(`\tclient: ${client}`);
    console.log(`\tservice provider: ${serviceProvider}`);

    await truffleDeployer.deploy(
      UptimeSLA,
      client,
      serviceProvider,
      LINK.address,
      Oracle.address,
      body.id,
      { value: 1000000000 });
  }).catch(console.log);
};
