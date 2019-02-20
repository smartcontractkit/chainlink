import cbor from 'cbor'
import abi from 'ethereumjs-abi'
import util from 'ethereumjs-util'
import BN from 'bn.js'

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
      buffer.length + 2
    )
  }

  return buffer
}

const bigNum = number => web3.utils.toBN(number)

const zeroX = value => (value.slice(0, 2) !== '0x' ? `0x${value}` : value)

const toHexWithoutPrefix = arg => {
  if (arg instanceof Buffer || arg instanceof BN) {
    return arg.toString('hex')
  } else if (arg instanceof Uint8Array) {
    return Array.prototype.reduce.call(
      arg,
      (a, v) => a + v.toString('16').padStart(2, '0'),
      ''
    )
  } else {
    return Buffer.from(arg, 'ascii').toString('hex')
  }
}

const toHex = value => {
  return zeroX(toHexWithoutPrefix(value))
}

const hexToAddress = hex => zeroX(bigNum(hex).toString('hex'))

export const assertActionThrows = action =>
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
        'VM Exception while processing transaction: revert'
      )
      assert(
        invalidOpcode || reverted,
        'expected following error message to include "invalid JUMP" or ' +
          `"revert": "${errorMessage}"`
      )
      // see https://github.com/ethereumjs/testrpc/issues/39
      // for why the "invalid JUMP" is the throw related error when using TestRPC
    })

export const decodeDietCBOR = data => {
  return cbor.decodeFirst(autoAddMapDelimiters(data))
}

export const decodeRunRequest = log => {
  const runABI = util.toBuffer(log.data)
  const types = [
    'address',
    'bytes32',
    'uint256',
    'address',
    'bytes4',
    'uint256',
    'uint256',
    'bytes'
  ]
  const [
    requester,
    requestId,
    payment,
    callbackAddress,
    callbackFunc,
    expiration,
    version,
    data
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
    data: autoAddMapDelimiters(data)
  }
}

export const functionSelector = signature =>
  '0x' +
  web3.utils
    .sha3(signature)
    .slice(2)
    .slice(0, 8)

export const fulfillOracleRequest = async (
  oracle,
  request,
  response,
  options
) => {
  if (!options) options = {}

  return oracle.fulfillOracleRequest(
    request.id,
    request.payment,
    request.callbackAddr,
    request.callbackFunc,
    request.expiration,
    response,
    options
  )
}

export const requestDataBytes = (specId, to, fHash, nonce, data) => {
  const types = [
    'address',
    'uint256',
    'bytes32',
    'address',
    'bytes4',
    'uint256',
    'uint256',
    'bytes'
  ]
  const values = [0, 0, specId, to, fHash, nonce, 1, data]
  const encoded = abiEncode(types, values)
  const funcSelector = functionSelector(
    'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)'
  )
  return funcSelector + encoded
}
