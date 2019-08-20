// truffle script

const request = require('request-promise').defaults({ jar: true })
const { deployer, abort, DEVNET_ADDRESS } = require('./common.js')
const {
  CHAINLINK_URL,
  ECHO_SERVER_URL
  // LINK_TOKEN_ADDRESS,
  // ORACLE_CONTRACT_ADDRESS
} = process.env
const { utils, web3 } = require('../evm/app/env.js')
const url = require('url')
const RunLog = artifacts.require('RunLog')
const LinkToken = artifacts.require('LinkToken')

process.env.SOLIDITY_INCLUDE = '../evm/contracts'

function futureOffset(seconds) {
  return parseInt(new Date().getTime() / 1000) + seconds
}

const main = async () => {
  const runLog = await RunLog.deployed()
  const linkToken = await LinkToken.deployed()

  const sessionsUrl = url.resolve(CHAINLINK_URL, '/sessions')
  const specsUrl = url.resolve(CHAINLINK_URL, '/v2/specs')
  const credentials = { email: 'notreal@fakeemail.ch', password: 'twochains' }

  // let RunLog = await deployer
  //   .perform(
  //     'contracts/RunLog.sol',
  //     LINK_TOKEN_ADDRESS,
  //     ORACLE_CONTRACT_ADDRESS
  //   )
  //   .catch(abort('Error deploying RunLog.sol'))
  // console.log(`Deployed RunLog at: ${RunLog.address}`)

  // const LinkToken = await deployer
  //   .load(
  //     '../node_modules/link_token/contracts/LinkToken.sol',
  //     LINK_TOKEN_ADDRESS
  //   )
  //   .catch(abort(`Error loading LinkToken at address ${LINK_TOKEN_ADDRESS}`))

  const amount = web3.utils.toBN(Number(utils.toWei('1000')).toString(16))
  await LinkToken.transfer(runLog.address, amount, {
    gas: 100000,
    from: DEVNET_ADDRESS
  }).catch(abort('Error transferring link to RunLog'))
  console.log(`Transferred ${amount} to RunLog at: ${runLog.address}`)

  const job = {
    _comment:
      'A runlog has a jobid baked into the contract so chainlink knows which job to run.',
    initiators: [{ type: 'runlog', params: { requesters: [runLog.address] } }],
    tasks: [
      // 10 seconds to ensure the time has not elapsed by the time the run is triggered
      { type: 'Sleep', params: { until: futureOffset(10) } },
      { type: 'HttpPost', params: { url: ECHO_SERVER_URL } },
      {
        type: 'EthTx',
        params: { functionSelector: 'fulfillOracleRequest(uint256,bytes32)' }
      }
    ]
  }
  await request.post(sessionsUrl, { json: credentials })
  let Job = await request
    .post(specsUrl, { json: job })
    .catch(abort('Error creating Job'))

  await runLog.request(web3.utils.asciiToHex(Job.data.id), {
    from: DEVNET_ADDRESS,
    gas: 2000000
  }).catch(abort('Error making RunLog request'))
  console.log(`Made RunLog request`)
}

module.exports = async callback => {
  try {
    await main()
    callback()
  } catch (error) {
    callback(error)
  }
}
