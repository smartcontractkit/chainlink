const clmigration = require('../clmigration.js')
const request = require('request-promise').defaults({ jar: true })
const UptimeSLA = artifacts.require('UptimeSLA')
const Oracle = artifacts.require('Oracle')
const LINK = artifacts.require('LinkToken')

const sessionsUrl = 'http://localhost:6688/sessions'
const specsUrl = 'http://localhost:6688/v2/specs'
const credentials = { email: 'notreal@fakeemail.ch', password: 'twochains' }
const job = {
  _comment:
    'GETs a number from JSON, multiplies by 10,000, and reports uint256',
  initiators: [{ type: 'runlog' }],
  tasks: [
    { type: 'httpGet' },
    { type: 'jsonParse' },
    { type: 'multiply', params: { times: 10000 } },
    { type: 'ethuint256' },
    { type: 'ethtx' },
  ],
}

module.exports = clmigration(async function(truffleDeployer) {
  const client = '0x542B68aE7029b7212A5223ec2867c6a94703BeE3'
  const serviceProvider = '0xB16E8460cCd76aEC437ca74891D3D358EA7d1d88'

  await request.post(sessionsUrl, { json: credentials })
  const body = await request.post(specsUrl, { json: job })
  console.log(`Deploying UptimeSLA:`)
  console.log(`\tjob: ${body.data.id}`)
  console.log(`\tclient: ${client}`)
  console.log(`\tservice provider: ${serviceProvider}`)

  await truffleDeployer.deploy(
    UptimeSLA,
    client,
    serviceProvider,
    LINK.address,
    Oracle.address,
    body.id,
    { value: 1000000000 },
  )
})
