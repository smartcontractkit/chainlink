// truffle script

const request = require('request-promise').defaults({ jar: true })
const url = require('url')
const { abort, DEVNET_ADDRESS, scriptRunner } = require('../common.js')
const RunLog = artifacts.require('RunLog')
const LinkToken = artifacts.require('LinkToken')
const { CHAINLINK_URL, ECHO_SERVER_URL } = process.env

const sessionsUrl = url.resolve(CHAINLINK_URL, '/sessions')
const specsUrl = url.resolve(CHAINLINK_URL, '/v2/specs')
const credentials = { email: 'notreal@fakeemail.ch', password: 'twochains' }
const amount = web3.utils.toBN(web3.utils.toWei('1000'))

const futureOffset = seconds => parseInt(new Date().getTime() / 1000) + seconds

const generateJob = requesterAddress => ({
  _comment:
    'A runlog has a jobid baked into the contract so chainlink knows which job to run.',
  initiators: [{ type: 'runlog', params: { requesters: [requesterAddress] } }],
  tasks: [
    // 10 seconds to ensure the time has not elapsed by the time the run is triggered
    { type: 'Sleep', params: { until: futureOffset(10) } },
    { type: 'HttpPost', params: { url: ECHO_SERVER_URL } },
    {
      type: 'EthTx',
      params: { functionSelector: 'fulfillOracleRequest(uint256,bytes32)' },
    },
  ],
})

const main = async () => {
  const runLog = await RunLog.deployed()
  const linkToken = await LinkToken.deployed()

  await linkToken
    .transfer(runLog.address, amount, {
      gas: 100000,
      from: DEVNET_ADDRESS,
    })
    .catch(abort('Error transferring link to RunLog'))
  console.log(`Transferred ${amount} to RunLog at: ${runLog.address}`)

  const job = generateJob(runLog.address)

  await request.post(sessionsUrl, { json: credentials })
  const Job = await request
    .post(specsUrl, { json: job })
    .catch(abort('Error creating Job'))

  await runLog
    .request(web3.utils.asciiToHex(Job.data.id), {
      from: DEVNET_ADDRESS,
      gas: 2000000,
    })
    .catch(abort('Error making RunLog request'))
  console.log(`Made RunLog request`)
}

module.exports = scriptRunner(main)
