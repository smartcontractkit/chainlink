import { eth } from '../../../../solidity/test/support/helpers'

// https://codepen.io/code_monk/pen/FvpfI
const randomHex = len => {
  const maxlen = 8
  const min = Math.pow(16, Math.min(len, maxlen) - 1)
  const max = Math.pow(16, Math.min(len, maxlen)) - 1
  const n = Math.floor(Math.random() * (max - min + 1)) + min
  let r = n.toString(16)
  while (r.length < len) {
    r = r + randomHex(len - maxlen)
  }
  return r
}

export const newAddress = () => ('0x' + randomHex(40))

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
