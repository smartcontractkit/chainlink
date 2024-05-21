/**
 * @packageDocumentation
 *
 * This file provides convenience functions to interact with existing solidity contract abstraction libraries, such as
 * @truffle/contract and ethers.js specifically for our `Oracle.sol` solidity smart contract.
 */
import { BigNumberish } from '@ethersproject/bignumber/lib/bignumber'
import { ethers } from 'ethers'
import { makeDebug } from './debug'
import { addCBORMapDelimiters, stripHexPrefix, toHex } from './helpers'
const debug = makeDebug('oracle')

/**
 * Transaction options such as gasLimit, gasPrice, data, ...
 */
type TxOptions = Omit<ethers.providers.TransactionRequest, 'to' | 'from'>

/**
 * A run request is an event emitted by `Oracle.sol` which triggers a job run
 * on a receiving chainlink node watching for RunRequests coming from that
 * specId + optionally requester.
 */
export interface RunRequest {
  /**
   * The ID of the job spec this request is targeting
   *
   * @solformat bytes32
   */
  specId: string
  /**
   * The requester of the run
   *
   * @solformat address
   */
  requester: string
  /**
   * The ID of the request, check Oracle.sol#oracleRequest to see how its computed
   *
   * @solformat bytes32
   */
  requestId: string
  /**
   * The amount of LINK used for payment
   *
   * @solformat uint256
   */
  payment: string
  /**
   * The address of the contract instance to callback with the fulfillment result
   *
   * @solformat address
   */
  callbackAddr: string
  /**
   * The function selector of the method that the oracle should call after fulfillment
   *
   * @solformat bytes4
   */
  callbackFunc: string
  /**
   * The expiration that the node should respond by before the requester can cancel
   *
   * @solformat uint256
   */
  expiration: string
  /**
   * The specified data version
   *
   * @solformat uint256
   */
  dataVersion: number
  /**
   * The CBOR encoded payload of the request
   *
   * @solformat bytes
   */
  data: Buffer

  /**
   * The hash of the signature of the OracleRequest event.
   * ```solidity
   *  event OracleRequest(
   *    bytes32 indexed specId,
   *    address requester,
   *    bytes32 requestId,
   *    uint256 payment,
   *    address callbackAddr,
   *    bytes4 callbackFunctionId,
   *    uint256 cancelExpiration,
   *    uint256 dataVersion,
   *    bytes data
   *  );
   * ```
   * Note: this is a property used for testing purposes only.
   * It is not part of the actual run request.
   *
   * @solformat bytes32
   */
  topic: string
}

/**
 * Convert the javascript format of the parameters needed to call the
 * ```solidity
 *  function fulfillOracleRequest(
 *    bytes32 _requestId,
 *    uint256 _payment,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _expiration,
 *    bytes32 _data
 *  )
 * ```
 * method on an Oracle.sol contract.
 *
 * @param runRequest The run request to flatten into the correct order to perform the `fulfillOracleRequest` function
 * @param response The response to fulfill the run request with, if it is an ascii string, it is converted to bytes32 string
 * @param txOpts Additional ethereum tx options
 */
export function convertFufillParams(
  runRequest: RunRequest,
  response: string,
  txOpts: TxOptions = {},
): [string, string, string, string, string, string, TxOptions] {
  const d = debug.extend('fulfillOracleRequestParams')
  d('Response param: %s', response)

  const bytes32Len = 32 * 2 + 2
  const convertedResponse =
    response.length < bytes32Len
      ? ethers.utils.formatBytes32String(response)
      : response
  d('Converted Response param: %s', convertedResponse)

  return [
    runRequest.requestId,
    runRequest.payment,
    runRequest.callbackAddr,
    runRequest.callbackFunc,
    runRequest.expiration,
    convertedResponse,
    txOpts,
  ]
}

/**
 * Convert the javascript format of the parameters needed to call the
 * ```solidity
 *  function fulfillOracleRequest2(
 *    bytes32 _requestId,
 *    uint256 _payment,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _expiration,
 *    bytes memory _data
 *  )
 * ```
 * method on an Oracle.sol contract.
 *
 * @param runRequest The run request to flatten into the correct order to perform the `fulfillOracleRequest` function
 * @param response The response to fulfill the run request with, if it is an ascii string, it is converted to bytes32 string
 * @param txOpts Additional ethereum tx options
 */
export function convertFulfill2Params(
  runRequest: RunRequest,
  responseTypes: string[],
  responseValues: string[],
  txOpts: TxOptions = {},
): [string, string, string, string, string, string, TxOptions] {
  const d = debug.extend('fulfillOracleRequestParams')
  d('Response param: %s', responseValues)
  const types = [...responseTypes]
  const values = [...responseValues]
  types.unshift('bytes32')
  values.unshift(runRequest.requestId)
  const convertedResponse = ethers.utils.defaultAbiCoder.encode(types, values)
  d('Encoded Response param: %s', convertedResponse)
  return [
    runRequest.requestId,
    runRequest.payment,
    runRequest.callbackAddr,
    runRequest.callbackFunc,
    runRequest.expiration,
    convertedResponse,
    txOpts,
  ]
}

/**
 * Convert the javascript format of the parameters needed to call the
 * ```solidity
 *  function cancelOracleRequest(
 *    bytes32 _requestId,
 *    uint256 _payment,
 *    bytes4 _callbackFunc,
 *    uint256 _expiration
 *  )
 * ```
 * method on an Oracle.sol contract.
 *
 * @param runRequest The run request to flatten into the correct order to perform the `cancelOracleRequest` function
 * @param txOpts Additional ethereum tx options
 */
export function convertCancelParams(
  runRequest: RunRequest,
  txOpts: TxOptions = {},
): [string, string, string, string, TxOptions] {
  return [
    runRequest.requestId,
    runRequest.payment,
    runRequest.callbackFunc,
    runRequest.expiration,
    txOpts,
  ]
}

/**
 * Convert the javascript format of the parameters needed to call the
 * ```solidity
 *  function cancelOracleRequestByRequester(
 *    uint256 nonce,
 *    uint256 _payment,
 *    bytes4 _callbackFunc,
 *    uint256 _expiration
 *  )
 * ```
 * method on an Oracle.sol contract.
 *
 * @param nonce The nonce used to generate the request ID
 * @param runRequest The run request to flatten into the correct order to perform the `cancelOracleRequest` function
 * @param txOpts Additional ethereum tx options
 */
export function convertCancelByRequesterParams(
  runRequest: RunRequest,
  nonce: number,
  txOpts: TxOptions = {},
): [number, string, string, string, TxOptions] {
  return [
    nonce,
    runRequest.payment,
    runRequest.callbackFunc,
    runRequest.expiration,
    txOpts,
  ]
}

/**
 * Abi encode parameters to call the `oracleRequest` method on the Oracle.sol contract.
 * ```solidity
 *  function oracleRequest(
 *    address _sender,
 *    uint256 _payment,
 *    bytes32 _specId,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _nonce,
 *    uint256 _dataVersion,
 *    bytes _data
 *  )
 * ```
 *
 * @param specId The Job Specification ID
 * @param callbackAddr The callback contract address for the response
 * @param callbackFunctionId The callback function id for the response
 * @param nonce The nonce sent by the requester
 * @param data The CBOR payload of the request
 */
export function encodeOracleRequest(
  specId: string,
  callbackAddr: string,
  callbackFunctionId: string,
  nonce: number,
  data: BigNumberish,
  dataVersion: BigNumberish = 1,
): string {
  const oracleRequestSighash = '0x40429946'
  return encodeRequest(
    oracleRequestSighash,
    specId,
    callbackAddr,
    callbackFunctionId,
    nonce,
    data,
    dataVersion,
  )
}

/**
 * Abi encode parameters to call the `operatorRequest` method on the Operator.sol contract.
 * ```solidity
 *  function operatorRequest(
 *    address _sender,
 *    uint256 _payment,
 *    bytes32 _specId,
 *    address _callbackAddress,
 *    bytes4 _callbackFunctionId,
 *    uint256 _nonce,
 *    uint256 _dataVersion,
 *    bytes _data
 *  )
 * ```
 *
 * @param specId The Job Specification ID
 * @param callbackAddr The callback contract address for the response
 * @param callbackFunctionId The callback function id for the response
 * @param nonce The nonce sent by the requester
 * @param data The CBOR payload of the request
 */
export function encodeRequestOracleData(
  specId: string,
  callbackFunctionId: string,
  nonce: number,
  data: BigNumberish,
  dataVersion: BigNumberish = 2,
): string {
  const sendOperatorRequestSigHash = '0x3c6d41b9'
  const requestInputs = [
    { name: '_sender', type: 'address' },
    { name: '_payment', type: 'uint256' },
    { name: '_specId', type: 'bytes32' },
    { name: '_callbackFunctionId', type: 'bytes4' },
    { name: '_nonce', type: 'uint256' },
    { name: '_dataVersion', type: 'uint256' },
    { name: '_data', type: 'bytes' },
  ]
  const encodedParams = ethers.utils.defaultAbiCoder.encode(
    requestInputs.map((i) => i.type),
    [
      ethers.constants.AddressZero,
      0,
      specId,
      callbackFunctionId,
      nonce,
      dataVersion,
      data,
    ],
  )
  return `${sendOperatorRequestSigHash}${stripHexPrefix(encodedParams)}`
}

function encodeRequest(
  oracleRequestSighash: string,
  specId: string,
  callbackAddr: string,
  callbackFunctionId: string,
  nonce: number,
  data: BigNumberish,
  dataVersion: BigNumberish = 1,
): string {
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
    oracleRequestInputs.map((i) => i.type),
    [
      ethers.constants.AddressZero,
      0,
      specId,
      callbackAddr,
      callbackFunctionId,
      nonce,
      dataVersion,
      data,
    ],
  )
  return `${oracleRequestSighash}${stripHexPrefix(encodedParams)}`
}

/**
 * Extract a javascript representation of a run request from the data
 * contained within a EVM log.
 * ```solidity
 *  event OracleRequest(
 *    bytes32 indexed specId,
 *    address requester,
 *    bytes32 requestId,
 *    uint256 payment,
 *    address callbackAddr,
 *    bytes4 callbackFunctionId,
 *    uint256 cancelExpiration,
 *    uint256 dataVersion,
 *    bytes data
 *  );
 * ```
 *
 * @param log The log to extract the run request from
 */
export function decodeRunRequest(log?: ethers.providers.Log): RunRequest {
  if (!log) {
    throw Error('No logs found to decode')
  }

  const ORACLE_REQUEST_TYPES = [
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
  ] = ethers.utils.defaultAbiCoder.decode(ORACLE_REQUEST_TYPES, log.data)

  return {
    specId: log.topics[1],
    requester,
    requestId: toHex(requestId),
    payment: toHex(payment),
    callbackAddr: callbackAddress,
    callbackFunc: toHex(callbackFunc),
    expiration: toHex(expiration),
    data: addCBORMapDelimiters(Buffer.from(stripHexPrefix(data), 'hex')),
    dataVersion: version.toNumber(),

    topic: log.topics[0],
  }
}

/**
 * Extract a javascript representation of a ConcreteChainlinked#Request event
 * from an EVM log.
 * ```solidity
 *  event Request(
 *    bytes32 id,
 *    address callbackAddress,
 *    bytes4 callbackfunctionSelector,
 *    bytes data
 *  );
 * ```
 * The request event is emitted from the `ConcreteChainlinked.sol` testing contract.
 *
 * @param log The log to decode
 */
export function decodeCCRequest(
  log: ethers.providers.Log,
): ethers.utils.Result {
  const d = debug.extend('decodeRunABI')
  d('params %o', log)

  const REQUEST_TYPES = ['bytes32', 'address', 'bytes4', 'bytes']
  const decodedValue = ethers.utils.defaultAbiCoder.decode(
    REQUEST_TYPES,
    log.data,
  )
  d('decoded value %o', decodedValue)

  return decodedValue
}
