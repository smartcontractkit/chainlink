const clmigration = require('../clmigration.js')
const request = require('request-promise').defaults({ jar: true })
const LinkToken = artifacts.require('LinkToken')
const Oracle = artifacts.require('Oracle')
const RunLog = artifacts.require('RunLog')

let sessionsUrl = 'http://localhost:6688/sessions'
let specsUrl = 'http://localhost:6688/v2/specs'
let credentials = { email: 'notreal@fakeemail.ch', password: 'twochains' }
let job = {
  _comment:
    'A runlog has a jobid baked into the contract so chainlink knows which job to run.',
  initiators: [{ type: 'runlog' }],
  tasks: [{ type: 'HttpPost', params: { url: 'http://localhost:6690' } }],
}

module.exports = clmigration(async function(truffleDeployer) {
  await request.post(sessionsUrl, { json: credentials })
  let body = await request.post(specsUrl, { json: job })
  console.log(`Deploying Consumer Contract with JobID ${body.data.id}`)
  let jobid = body.data.id
  if (jobid && jobid.slice(0, 2) != '0x') {
    jobid = `0x${jobid}` // hack to prefix 0x to satisfy bytes32 requirement
  }
  await truffleDeployer.deploy(RunLog, LinkToken.address, Oracle.address, jobid)
})
