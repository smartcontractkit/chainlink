let chainlinkDeployer = require("../chainlink_deployer.js");
let Consumer = artifacts.require("./Consumer.sol");
let Oracle = artifacts.require("./Oracle.sol");
let LinkToken = artifacts.require("./LinkToken.sol");
let fs = require('fs');

let url = "http://chainlink:twochains@localhost:6688/v2/specs";
let job = {
  "initiators": [{ "type": "runlog" }],
  "tasks": [
    { "type": "httpGet" },
    { "type": "jsonParse" },
    { "type": "multiply", "times": 100 },
    { "type": "ethuint256" },
    { "type": "ethtx" }
  ]
}

module.exports = function(truffleDeployer) {
  truffleDeployer.then(async () => {
    let body = await chainlinkDeployer.job(url, job);
    console.log(`Deploying Consumer:`);
    console.log(`\tjob: ${body.id}`);
    await truffleDeployer.deploy(Consumer, LinkToken.address, Oracle.address, body.id);
  }).catch(console.log);
};
