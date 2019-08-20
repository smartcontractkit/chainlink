// truffle script

const request = require('request-promise').defaults({ jar: true })
const url = require('url')
const { abort, DEVNET_ADDRESS } = require('./common.js')
const EthLog = artifacts.require('EthLog')
const { CHAINLINK_URL, ECHO_SERVER_URL } = process.env

const main = async () => {
  const ethLog = await EthLog.deployed()

  const sessionsUrl = url.resolve(CHAINLINK_URL, '/sessions')
  const credentials = { email: 'notreal@fakeemail.ch', password: 'twochains' }
  await request.post(sessionsUrl, { json: credentials })

  const job = {
    _comment: 'An ethlog with no address listens to all addresses.',
    initiators: [{ type: 'ethlog', params: { address: ethLog.address } }],
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

// truffle exec won't capture errors automatically
module.exports = async callback => {
  try {
    await main()
    callback()
  } catch (error) {
    callback(error)
  }
}
