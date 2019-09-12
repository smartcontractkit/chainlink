/* eslint-disable @typescript-eslint/no-var-requires */

const bigNum = number => web3.utils.toBN(number)
const cbor = require('cbor')
const abi = require('ethereumjs-abi')
const util = require('ethereumjs-util')
const BN = require('bn.js')
const ethJSUtils = require('ethereumjs-util')

web3.providers.HttpProvider.prototype.sendAsync =
  web3.providers.HttpProvider.prototype.send

const exps = {}

const sendEth = (method, params) =>
  new Promise((resolve, reject) => {
    web3.currentProvider.sendAsync(
      {
        jsonrpc: '2.0',
        method: method,
        params: params || [],
        id: new Date().getTime(),
      },
      (error, response) => (error ? reject(error) : resolve(response.result)),
      () => {},
      () => {},
    )
  })

exps.assertBigNum = (a, b, failureMessage) =>
  assert(
    bigNum(a).eq(bigNum(b)),
    `BigNum ${a} is not ${b}` + (failureMessage ? ': ' + failureMessage : ''),
  )

exps.getLatestTimestamp = async () => {
  const latestBlock = await web3.eth.getBlock('latest', false)
  return web3.utils.toDecimal(latestBlock.timestamp)
}

exps.fastForwardTo = async target => {
  const now = await exps.getLatestTimestamp()
  assert.isAbove(target, now, 'Cannot fast forward to the past')
  const difference = target - now
  await sendEth('evm_increaseTime', [difference])
  await sendEth('evm_mine')
}

const minutes = number => number * 60
const hours = number => number * minutes(60)

exps.days = number => number * hours(24)

const abiEncode = (types, values) => {
  return abi.rawEncode(types, values).toString('hex')
}

const startMapBuffer = Buffer.from([0xbf])
const endMapBuffer = Buffer.from([0xff])

const autoAddMapDelimiters = data => {
  let buffer = data

  if (buffer[0] >> 5 !== 5) {
    buffer = Buffer.concat(
      [startMapBuffer, buffer, endMapBuffer],
      buffer.length + 2,
    )
  }

  return buffer
}

const zeroX = value => (value.slice(0, 2) !== '0x' ? `0x${value}` : value)

const toHexWithoutPrefix = arg => {
  if (arg instanceof Buffer || arg instanceof BN) {
    return arg.toString('hex')
  } else if (arg instanceof Uint8Array) {
    return Array.prototype.reduce.call(
      arg,
      (a, v) => a + v.toString('16').padStart(2, '0'),
      '',
    )
  } else {
    return Buffer.from(arg, 'ascii').toString('hex')
  }
}

const toHex = value => {
  return zeroX(toHexWithoutPrefix(value))
}

exps.assertActionThrows = action =>
  Promise.resolve()
    .then(action)
    .catch(error => {
      assert(error, 'Expected an error to be raised')
      assert(error.message, 'Expected an error to be raised')
      return error.message
    })
    .then(errorMessage => {
      assert(errorMessage, 'Expected an error to be raised')
      const invalidOpcode = errorMessage.includes('invalid opcode')
      const reverted = errorMessage.includes(
        'VM Exception while processing transaction: revert',
      )
      assert(
        invalidOpcode || reverted,
        'expected following error message to include "invalid JUMP" or ' +
          `"revert": "${errorMessage}"`,
      )
      // see https://github.com/ethereumjs/testrpc/issues/39
      // for why the "invalid JUMP" is the throw related error when using TestRPC
    })

exps.decodeDietCBOR = data => {
  return cbor.decodeFirst(autoAddMapDelimiters(ethJSUtils.toBuffer(data)))
}

exps.decodeRunRequest = log => {
  const runABI = util.toBuffer(log.data)
  const types = [
    'address',
    'bytes32',
    'uint256',
    'address',
    'bytes4',
    'uint256',
    'uint256',
    'bytes',
  ]
  const [
    requester,
    requestId,
    payment,
    callbackAddress,
    callbackFunc,
    expiration,
    version,
    data,
  ] = abi.rawDecode(types, runABI)

  return {
    topic: log.topics[0],
    jobId: log.topics[1],
    requester: zeroX(requester),
    id: toHex(requestId),
    payment: toHex(payment),
    callbackAddr: zeroX(callbackAddress),
    callbackFunc: toHex(callbackFunc),
    expiration: toHex(expiration),
    dataVersion: version,
    data: autoAddMapDelimiters(data),
  }
}

exps.functionSelector = signature =>
  '0x' +
  web3.utils
    .sha3(signature)
    .slice(2)
    .slice(0, 8)

exps.fulfillOracleRequest = async (oracle, request, response, options) => {
  if (!options) options = {}

  return oracle.fulfillOracleRequest(
    request.id,
    request.payment,
    request.callbackAddr,
    request.callbackFunc,
    request.expiration,
    response,
    options,
  )
}

exps.requestDataBytes = (specId, to, fHash, nonce, data) => {
  const types = [
    'address',
    'uint256',
    'bytes32',
    'address',
    'bytes4',
    'uint256',
    'uint256',
    'bytes',
  ]
  const values = [0, 0, specId, to, fHash, nonce, 1, data]
  const encoded = abiEncode(types, values)
  const funcSelector = exps.functionSelector(
    'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)',
  )
  return funcSelector + encoded
}

module.exports = exps
