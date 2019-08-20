// truffle script

const request = require('request-promise').defaults({ jar: true })
// const { deployer, abort, DEVNET_ADDRESS } = require('./common.js')
const { abort, DEVNET_ADDRESS } = require('./common.js')
const { CHAINLINK_URL, ECHO_SERVER_URL, ETH_LOG_ADDRESS } = process.env
const url = require('url')
const EthLog = artifacts.require('EthLog')

const main = async () => {
  // let EthLog = await deployer
  //   .perform('contracts/EthLog.sol')
  //   .catch(abort('Error deploying EthLog.sol'))
  // console.log(`Deployed EthLog at: ${EthLog.address}`)
  const ethLog = await EthLog.deployed()

  const sessionsUrl = url.resolve(CHAINLINK_URL, '/sessions')
  const credentials = { email: 'notreal@fakeemail.ch', password: 'twochains' }
  await request.post(sessionsUrl, { json: credentials })

  const job = {
    _comment: 'An ethlog with no address listens to all addresses.',
    initiators: [{ type: 'ethlog', params: { address: ETH_LOG_ADDRESS } }],
    tasks: [{ type: 'HttpPost', params: { url: ECHO_SERVER_URL } }]
  }
  const specsUrl = url.resolve(CHAINLINK_URL, '/v2/specs')
  let Job = await request
  .post(specsUrl, { json: job })
    .catch(abort('Error creating Job'))

  console.log('Deployed Job at:', Job.data.id)

  await ethLog
    .logEvent({ from: DEVNET_ADDRESS, gas: 200000 })
    .catch(abort('Error making EthLog entry'))
  console.log(`Made EthLog entry`)
}

module.exports = async callback => {
  try {
    await main()
    callback()
  } catch (error) {
    callback(error)
  }
}
