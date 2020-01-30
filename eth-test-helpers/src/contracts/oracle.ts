import { ContractTransaction, ethers } from 'ethers'
import { makeDebug } from '../debug'
import { addCBORMapDelimiters, stripHexPrefix, toHex } from '../helpers'

const debug = makeDebug('oracle')
type TxOptions = Omit<ethers.providers.TransactionRequest, 'to' | 'from'>

// TODO find ethers equivalent
class TransactionOverrides {
  nonce?: ethers.utils.BigNumberish | Promise<ethers.utils.BigNumberish>
  gasLimit?: ethers.utils.BigNumberish | Promise<ethers.utils.BigNumberish>
  gasPrice?: ethers.utils.BigNumberish | Promise<ethers.utils.BigNumberish>
  value?: ethers.utils.BigNumberish | Promise<ethers.utils.BigNumberish>
  chainId?: number | Promise<number>
}

export interface RunRequest {
  callbackAddr: string
  callbackFunc: string
  data: Buffer
  dataVersion: number
  expiration: string
  id: string
  jobId: string
  payment: string
  requester: string
  topic: string
}

interface Fulfillable {
  fulfillOracleRequest(
    _requestId: ethers.utils.Arrayish,
    _payment: ethers.utils.BigNumberish,
    _callbackAddress: string,
    _callbackFunctionId: ethers.utils.Arrayish,
    _expiration: ethers.utils.BigNumberish,
    _data: ethers.utils.Arrayish,
    overrides?: TransactionOverrides,
  ): Promise<ContractTransaction>
}
export async function fulfillOracleRequest(
  oracleContract: Fulfillable,
  runRequest: RunRequest,
  response: string,
  options: TxOptions = {
    gasLimit: 1000000, // FIXME: incorrect gas estimation
  },
): ReturnType<typeof oracleContract.fulfillOracleRequest> {
  const d = debug.extend('fulfillOracleRequest')
  d('Response param: %s', response)

  const bytes32Len = 32 * 2 + 2
  const convertedResponse =
    response.length < bytes32Len
      ? ethers.utils.formatBytes32String(response)
      : response
  d('Converted Response param: %s', convertedResponse)

  return oracleContract.fulfillOracleRequest(
    runRequest.id,
    runRequest.payment,
    runRequest.callbackAddr,
    runRequest.callbackFunc,
    runRequest.expiration,
    convertedResponse,
    options,
  )
}

interface Cancellable {
  cancelOracleRequest(
    _requestId: ethers.utils.Arrayish,
    _payment: ethers.utils.BigNumberish,
    _callbackFunc: ethers.utils.Arrayish,
    _expiration: ethers.utils.BigNumberish,
    overrides?: TransactionOverrides,
  ): Promise<ContractTransaction>
}
export async function cancelOracleRequest(
  oracleContract: Cancellable,
  request: RunRequest,
  options: TxOptions = {},
): ReturnType<typeof oracleContract.cancelOracleRequest> {
  return oracleContract.cancelOracleRequest(
    request.id,
    request.payment,
    request.callbackFunc,
    request.expiration,
    options,
  )
}

export function encodeOracleRequest(
  specId: string,
  to: string,
  fHash: string,
  nonce: number,
  dataBytes: string,
): string {
  // 'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)'
  const oracleRequestSighash = '0x40429946'
  const oracleRequestInputs = [
    { name: '_sender', type: 'address' },
    { name: '_payment', type: 'uint256' },
    { name: '_specId', type: 'bytes32' },
    { name: '_callbackAddress', type: 'address' },
    { name: '_callbackFunctionId', type: 'bytes4' },
    { name: '_nonce', type: 'uint256' },
    { name: '_dataVersion', type: 'uint256' },
    { name: '_data', type: 'bytes' },
  ]

  const encodedParams = ethers.utils.defaultAbiCoder.encode(
    oracleRequestInputs.map(i => i.type),
    [ethers.constants.AddressZero, 0, specId, to, fHash, nonce, 1, dataBytes],
  )

  return `${oracleRequestSighash}${stripHexPrefix(encodedParams)}`
}

export function decodeRunRequest(log?: ethers.providers.Log): RunRequest {
  if (!log) {
    throw Error('No logs found to decode')
  }

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
  ] = ethers.utils.defaultAbiCoder.decode(types, log.data)

  return {
    callbackAddr: callbackAddress,
    callbackFunc: toHex(callbackFunc),
    data: addCBORMapDelimiters(Buffer.from(stripHexPrefix(data), 'hex')),
    dataVersion: version.toNumber(),
    expiration: toHex(expiration),
    id: toHex(requestId),
    jobId: log.topics[1],
    payment: toHex(payment),
    requester,
    topic: log.topics[0],
  }
}

/**
 * Decode a log into a run
 * @param log The log to decode
 * @todo Do we really need this?
 */
export function decodeRunABI(
  log: ethers.providers.Log,
): [string, string, string, string] {
  const d = debug.extend('decodeRunABI')
  d('params %o', log)

  const types = ['bytes32', 'address', 'bytes4', 'bytes']
  const decodedValue = ethers.utils.defaultAbiCoder.decode(types, log.data)
  d('decoded value %o', decodedValue)

  return decodedValue
}
