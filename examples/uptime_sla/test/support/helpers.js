module.exports = {}

web3.providers.HttpProvider.prototype.sendAsync =
  web3.providers.HttpProvider.prototype.send

const getLatestTimestamp = async () => {
  const latestBlock = await web3.eth.getBlock('latest', false)
  return web3.utils.toDecimal(latestBlock.timestamp)
}

module.exports.getLatestTimestamp = getLatestTimestamp

const sendEth = (method, params) =>
  new Promise((resolve, reject) => {
    web3.currentProvider.sendAsync(
      {
        jsonrpc: '2.0',
        method: method,
        params: params || [],
        id: new Date().getTime()
      },
      (error, response) => (error ? reject(error) : resolve(response.result)),
      () => {},
      () => {}
    )
  })

const fastForwardTo = async target => {
  const now = await getLatestTimestamp()
  assert.isAbove(target, now, 'Cannot fast forward to the past')
  const difference = target - now
  await sendEth('evm_increaseTime', [difference])
  await sendEth('evm_mine')
}

module.exports.fastForwardTo = fastForwardTo

const minutes = number => number * 60
const hours = number => number * minutes(60)

const days = number => number * hours(24)

module.exports.days = days
