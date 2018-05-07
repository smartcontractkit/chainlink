const request = require("request-promise");
const LinkToken = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/LinkToken.sol");
const Oracle = artifacts.require("../node_modules/smartcontractkit/chainlink/solidity/contracts/Oracle.sol");
const RunLog = artifacts.require("./RunLog.sol");

let url = "http://chainlink:twochains@localhost:6688/v2/specs";
let job = {
  "_comment": "A runlog has a jobid baked into the contract so chainlink knows which job to run.",
  "initiators": [{ "type": "runlog" }],
  "tasks": [
    { "type": "HttpPost", "url": "http://localhost:6690" }
  ]
};

module.exports = function(truffleDeployer) {
  truffleDeployer.then(async () => {
    let body = await request.post(url, {json: job});
    console.log(`Deploying Consumer Contract with JobID ${body.id}`);
    await truffleDeployer.deploy(RunLog, LinkToken.address, Oracle.address, body.id);
  }).catch(console.log);
};
