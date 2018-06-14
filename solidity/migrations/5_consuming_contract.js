const clmigration = require('../clmigration.js')
const request = require('request-promise')
const Consumer = artifacts.require('./Consumer.sol')
const Oracle = artifacts.require('./Oracle.sol')
const LinkToken = artifacts.require('./LinkToken.sol')

let url = 'http://chainlink:twochains@localhost:6688/v2/specs'
let job = {
  'initiators': [{ 'type': 'runlog' }],
  'tasks': [
    { 'type': 'httpGet' },
    { 'type': 'jsonParse' },
    { 'type': 'multiply', 'times': 100 },
    { 'type': 'ethuint256' },
    { 'type': 'ethtx' }
  ]
}

module.exports = clmigration(async function (truffleDeployer) {
  let body = await request.post(url, {json: job})
  console.log(`Deploying Consumer:`)
  console.log(`\tjob: ${body.id}`)
  await truffleDeployer.deploy(Consumer, LinkToken.address, Oracle.address, body.id)
})
