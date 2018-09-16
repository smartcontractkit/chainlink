import { eth } from '../../../../solidity/test/support/helpers'

export const getLatestTimestamp = async () => {
  const latestBlock = await eth.getBlock('latest', false)
  return web3.toDecimal(latestBlock.timestamp)
}

const sendEth = (method, params) => (
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
)

export const fastForwardTo = async target => {
  const now = await getLatestTimestamp()
  assert.isAbove(target, now, 'Cannot fast forward to the past')
  const difference = target - now
  await sendEth('evm_increaseTime', [difference])
  await sendEth('evm_mine')
}

const minutes = number => number * 60
const hours = number => (number * minutes(60))
export const days = number => (number * hours(24))
