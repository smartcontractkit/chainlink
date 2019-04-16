import cbor from 'cbor'
import { join, resolve as pathResolve } from 'path'
import { assertBigNum } from './matchers'

const contractPathHead = pathResolve(join(__dirname, '/../..'))

// Paths for finding solidity files during compilation via deployer.
// See compile.js for more info.
process.env.SOLIDITY_INCLUDE = [
  'contracts/',
  'contracts/examples/',
  'contracts/interfaces/',
  '../node_modules/',
  '../node_modules/link_token/contracts',
  '../node_modules/openzeppelin-solidity/contracts/ownership/',
  '../node_modules/@ensdomains/ens/contracts/',
]
  .map(p => join(contractPathHead, p))
  .join(':')

// Relative paths needed for chainlink/examples/{uptime_sla,echo_server}
// Note that these will be relative to the truffle script
// (usually executed as ./node_modules/.bin/truffle)
process.env.SOLIDITY_INCLUDE += ':../../../../:../../contracts'

// Key for defaultAccount, defined below.
const PRIVATE_KEY =
  'c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3'

/* tslint:disable no-var-requires */
const Wallet = require('../../app/wallet.js')
const Utils = require('../../app/utils.js')
const Deployer = require('../../app/deployer.js')

const abi = require('ethereumjs-abi')
const util = require('ethereumjs-util')
const BN = require('bn.js')
const ethjsUtils = require('ethereumjs-util')
/* tslint:enable no-var-requires */

const HEX_BASE = 16

// https://github.com/ethereum/web3.js/issues/1119#issuecomment-394217563
web3.providers.HttpProvider.prototype.sendAsync =
  web3.providers.HttpProvider.prototype.send
export const eth = web3.eth

const INVALIDVALUE = {
  // If you got this value, you probably tried to use one of the variables below
  // before they were initialized. Do any test initialization which requires
  // them in a callback passed to Mocha's `before` or `beforeEach`.
  // https://mochajs.org/#asynchronous-hooks
  unitializedValueProbablyShouldUseVaribleInMochaBeforeCallback: null
}

export let [
  accounts,
  defaultAccount,
  oracleNode1,
  oracleNode2,
  oracleNode3,
  stranger,
  consumer,
  oracleNode
] = Array(1000).fill(INVALIDVALUE)

before(async function queryEthClientForConstants() {
  accounts = await eth.getAccounts()
  ;[
    defaultAccount,
    oracleNode1,
    oracleNode2,
    oracleNode3,
    stranger,
    consumer
  ] = accounts.slice(0, 6)
  oracleNode = oracleNode1
})

export const utils = Utils(web3.currentProvider)
export const wallet = Wallet(PRIVATE_KEY, utils)
export const deployer = Deployer(wallet, utils)

const bNToStringOrIdentity = (a: any): any => (BN.isBN(a) ? a.toString() : a)

// Deal with transfer amount type truffle doesn't currently handle. (BN)
export const wrappedERC20 = (contract: any): any => ({
  ...contract,
  transfer: async (address: any, amount: any) =>
    contract.transfer(address, bNToStringOrIdentity(amount)),
  transferAndCall: async (
    address: any,
    amount: any,
    payload: any,
    options: any
  ) =>
    contract.transferAndCall(
      address,
      bNToStringOrIdentity(amount),
      payload,
      options
    )
})

export const linkContract = async (): Promise<any> => {
  return wrappedERC20(await deploy('link_token/contracts/LinkToken.sol'))
}

export const bigNum = (num: any): BigNumber => web3.utils.toBN(num)
assertBigNum(
  bigNum('1'),
  bigNum(1),
  'Different representations should give same BNs'
)

// toWei(n) is n * 10**18, as a BN.
export const toWei = (num: string | number): any =>
  bigNum(web3.utils.toWei(bigNum(num)))
assertBigNum(
  toWei('1'),
  toWei(1),
  'Different representations should give same BNs'
)

export const toUtf8 = web3.utils.toUtf8

export const keccak = web3.utils.sha3

export const hexToInt = (str: string): any => bigNum(str).toNumber()

export const toHexWithoutPrefix = (arg: any): string => {
  if (arg instanceof Buffer || arg instanceof BN) {
    return arg.toString('hex')
  } else if (arg instanceof Uint8Array) {
    return Array.prototype.reduce.call(
      arg,
      (a: any, v: any) => a + v.toString('16').padStart(2, '0'),
      ''
    )
  } else if (Number(arg) === arg) {
    return arg.toString(16).padStart(64, '0')
  } else {
    return Buffer.from(arg, 'ascii').toString('hex')
  }
}

export const toHex = (value: any): string => {
  return Ox(toHexWithoutPrefix(value))
}

export const Ox = (value: any): string =>
  value.slice(0, 2) !== '0x' ? `0x${value}` : value

// True if h is a standard representation of a byte array, false otherwise
export const isByteRepresentation = (h: any): boolean => {
  return h instanceof Buffer || h instanceof BN || h instanceof Uint8Array
}

export const deploy = async (filePath: any, ...args: any[]) =>
  deployer.perform(filePath, ...args)

export const getEvents = (contract: any): Promise<any[]> =>
  new Promise((resolve, reject) =>
    contract
      .allEvents()
      .get((error: any, events: any) =>
        error ? reject(error) : resolve(events)
      )
  )

export const getLatestEvent = async (contract: any): Promise<any[]> => {
  const events = await getEvents(contract)
  return events[events.length - 1]
}

// link param must be from linkContract(), if amount is a BN
export const requestDataFrom = (
  oc: any,
  link: any,
  amount: any,
  args: any,
  options: any
): any => {
  if (!options) {
    options = {}
  }
  return link.transferAndCall(oc.address, amount, args, options)
}

export const functionSelector = (signature: any): string =>
  '0x' +
  keccak(signature)
    .slice(2)
    .slice(0, 8)

export const assertActionThrows = (action: any) =>
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

export const checkPublicABI = (contract: any, expectedPublic: any) => {
  const actualPublic = []
  for (const method of contract.abi) {
    if (method.type === 'function') {
      actualPublic.push(method.name)
    }
  }

  for (const method of actualPublic) {
    const index = expectedPublic.indexOf(method)
    assert.isAtLeast(index, 0, `#${method} is NOT expected to be public`)
  }

  for (const method of expectedPublic) {
    const index = actualPublic.indexOf(method)
    assert.isAtLeast(index, 0, `#${method} is expected to be public`)
  }
}

export const decodeRunABI = (log: any): any => {
  const runABI = util.toBuffer(log.data)
  const types = ['bytes32', 'address', 'bytes4', 'bytes']
  return abi.rawDecode(types, runABI)
}

const startMapBuffer = Buffer.from([0xbf])
const endMapBuffer = Buffer.from([0xff])

export const decodeRunRequest = (log: any): any => {
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
    callbackAddr: Ox(callbackAddress),
    callbackFunc: toHex(callbackFunc),
    data: autoAddMapDelimiters(data),
    dataVersion: version,
    expiration: toHex(expiration),
    id: toHex(requestId),
    jobId: log.topics[1],
    payment: toHex(payment),
    requester: Ox(requester),
    topic: log.topics[0]
  }
}

const autoAddMapDelimiters = (data: any): Buffer => {
  let buffer = data

  if (buffer[0] >> 5 !== 5) {
    buffer = Buffer.concat(
      [startMapBuffer, buffer, endMapBuffer],
      buffer.length + 2
    )
  }

  return buffer
}

export const decodeDietCBOR = (data: any): any => {
  return cbor.decodeFirstSync(autoAddMapDelimiters(data))
}

export const runRequestId = (log: any): any => {
  const { requestId } = decodeRunRequest(log)
  return requestId
}

export const requestDataBytes = (
  specId: any,
  to: any,
  fHash: any,
  nonce: any,
  data: any
): any => {
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

export const abiEncode = (types: any, values: any): string => {
  return abi.rawEncode(types, values).toString('hex')
}

export const newUint8ArrayFromStr = (str: string): Uint8Array => {
  const codePoints = Array.prototype.map.call(str, (c: string) =>
    c.charCodeAt(0)
  )
  return Uint8Array.from(codePoints)
}

// newUint8Array returns a uint8array of count bytes from either a hex or
// decimal string, hex strings must begin with 0x
export const newUint8Array = (str: string, count: number): any => {
  let result = new Uint8Array(count)

  if (str.startsWith('0x') || str.startsWith('0X')) {
    const hexStr = str.slice(2).padStart(count * 2, '0')
    for (let i = result.length; i >= 0; i--) {
      const offset = i * 2
      result[i] = parseInt(hexStr[offset] + hexStr[offset + 1], HEX_BASE)
    }
  } else {
    const num = bigNum(str)
    result = newHash('0x' + num.toString(HEX_BASE))
  }

  return result
}

// newSignature returns a signature object with v, r, and s broken up
export const newSignature = (str: string): any => {
  const oracleSignature = newUint8Array(str, 65)
  let v = oracleSignature[64]
  if (v < 27) {
    v += 27
  }
  return {
    full: oracleSignature,
    r: oracleSignature.slice(0, 32),
    s: oracleSignature.slice(32, 64),
    v
  }
}

// newHash returns a 65 byte Uint8Array for representing a hash
export const newHash = (str: string): Uint8Array => {
  return newUint8Array(str, 32)
}

// newAddress returns a 20 byte Uint8Array for representing an address
export const newAddress = (str: string): Uint8Array => {
  return newUint8Array(str, 20)
}

// lengthTypedArrays sums the length of all specified TypedArrays
export const lengthTypedArrays = <T>(
  ...arrays: Array<ArrayLike<T>>
): number => {
  return arrays.reduce((a, v) => a + v.length, 0)
}

export const toBuffer = (uint8a: Uint8Array): Buffer => {
  return Buffer.from(uint8a)
}

// concatTypedArrays recursively concatenates TypedArrays into one big
// TypedArray
// TODO: Does not work recursively
export const concatTypedArrays = <T>(
  ...arrays: Array<ArrayLike<T>>
): ArrayLike<T> => {
  const size = lengthTypedArrays(...arrays)
  const arrayCtor: any = arrays[0].constructor
  const result = new arrayCtor(size)
  let offset = 0
  arrays.forEach(a => {
    result.set(a, offset)
    offset += a.length
  })
  return result
}

export const increaseTime5Minutes = async () => {
  await web3.currentProvider.send(
    {
      id: 0,
      jsonrpc: '2.0',
      method: 'evm_increaseTime',
      params: [300]
    },
    (error: any, result: any) => {
      if (error) {
        // tslint:disable-next-line:no-console
        console.log(`Error during helpers.increaseTime5Minutes! ${error}`)
        throw error
      }
    }
  )
}

export const sendToEvm = async (evmMethod: string, ...params: any) => {
  await web3.currentProvider.sendAsync(
    {
      id: 0,
      jsonrpc: '2.0',
      method: evmMethod,
      params: [...params]
    },
    (error: any, result: any) => {
      if (error) {
        // tslint:disable-next-line:no-console
        console.log(`Error during ${evmMethod}! ${error}`)
        throw error
      }
    }
  )
}

export const mineBlocks = async (blocks: number) => {
  for (let i = 0; i < blocks; i++) {
    await sendToEvm('evm_mine')
  }
}

export const createTxData = (
  selector: string,
  types: any,
  values: any
): any => {
  const funcSelector = functionSelector(selector)
  const encoded = abiEncode([...types], [...values])
  return funcSelector + encoded
}

export const calculateSAID = ({
  payment,
  expiration,
  endAt,
  oracles,
  requestDigest
}: any): Uint8Array => {
  const serviceAgreementIDInput = concatTypedArrays(
    newHash(payment.toString()),
    newHash(expiration.toString()),
    newHash(endAt.toString()),
    concatTypedArrays(
      ...oracles
        .map(newAddress)
        .map(toHex)
        .map(newHash)
    ),
    requestDigest
  )
  const serviceAgreementIDInputDigest = ethjsUtils.keccak(
    toHex(serviceAgreementIDInput)
  )
  return newHash(toHex(serviceAgreementIDInputDigest))
}

export const recoverPersonalSignature = (
  message: Uint8Array,
  signature: any
): any => {
  const personalSignPrefix = newUint8ArrayFromStr(
    '\x19Ethereum Signed Message:\n'
  )
  const personalSignMessage = Uint8Array.from(
    concatTypedArrays(
      personalSignPrefix,
      newUint8ArrayFromStr(message.length.toString()),
      message
    )
  )
  const digest = ethjsUtils.keccak(toBuffer(personalSignMessage))
  const requestDigestPubKey = ethjsUtils.ecrecover(
    digest,
    signature.v,
    toBuffer(signature.r),
    toBuffer(signature.s)
  )
  return ethjsUtils.pubToAddress(requestDigestPubKey)
}

export const personalSign = async (
  account: any,
  message: any
): Promise<any> => {
  if (!isByteRepresentation(message)) {
    throw new Error(`Message ${message} is not a recognized representation of a byte array.
    (Can be Buffer, BigNumber, Uint8Array, 0x-prepended hexadecimal string.)`)
  }
  return newSignature(await web3.eth.sign(toHex(message), account))
}

export const executeServiceAgreementBytes = (
  sAID: any,
  to: any,
  fHash: any,
  nonce: any,
  data: any
): any => {
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
  const values = [0, 0, sAID, to, fHash, nonce, 1, data]
  const encoded = abiEncode(types, values)
  const funcSelector = functionSelector(
    'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)'
  )
  return funcSelector + encoded
}

// Convenience functions for constructing hexadecimal representations of
// binary serializations.
export const padHexTo256Bit = (s: string): string => s.padStart(64, '0')
export const strip0x = (s: string): string =>
  s.startsWith('0x') ? s.slice(2) : s
export const pad0xHexTo256Bit = (s: string): string =>
  padHexTo256Bit(strip0x(s))
export const padNumTo256Bit = (n: number): string =>
  padHexTo256Bit(n.toString(16))

export const initiateServiceAgreementArgs = ({
  payment,
  expiration,
  endAt,
  oracles,
  oracleSignatures,
  requestDigest
}: any): any[] => [
  toHex(newHash(payment.toString())),
  toHex(newHash(expiration.toString())),
  toHex(newHash(endAt.toString())),
  oracles.map(newAddress).map(toHex),
  oracleSignatures.map((os: any) => os.v),
  oracleSignatures.map((os: any) => toHex(os.r)),
  oracleSignatures.map((os: any) => toHex(os.s)),
  toHex(requestDigest)
]

// Call coordinator contract to initiate the specified service agreement, and
// get the return value
export const initiateServiceAgreementCall = async (
  coordinator: any,
  args: any
): Promise<any> =>
  coordinator.initiateServiceAgreement.call(
    ...initiateServiceAgreementArgs(args)
  )

/** Call coordinator contract to initiate the specified service agreement. */
export const initiateServiceAgreement = async (
  coordinator: any,
  args: any
): Promise<any> =>
  coordinator.initiateServiceAgreement(...initiateServiceAgreementArgs(args))

/** Check that the given service agreement was stored at the correct location */
export const checkServiceAgreementPresent = async (
  coordinator: any,
  { payment, expiration, endAt, requestDigest, id }: any
): Promise<any> => {
  const sa = await coordinator.serviceAgreements.call(id)
  assertBigNum(sa[0], bigNum(payment))
  assertBigNum(sa[1], bigNum(expiration))
  assertBigNum(sa[2], bigNum(endAt))
  assert.equal(sa[3], toHex(requestDigest))

  /// / TODO:

  /// / Web3.js doesn't support generating an artifact for arrays
  /// within a struct. / This means that we aren't returned the
  /// list of oracles and / can't assert on their values.
  /// /

  /// / However, we can pass them into the function to generate the
  /// ID / & solidity won't compile unless we pass the correct
  /// number and / type of params when initializing the
  /// ServiceAgreement struct, / so we have some indirect test
  /// coverage.
  /// /
  /// / https://github.com/ethereum/web3.js/issues/1241
  /// / assert.equal(
  /// /   sa[2],
  /// /   ['0x70AEc4B9CFFA7b55C0711b82DD719049d615E21d',
  /// /    '0xd26114cd6EE289AccF82350c8d8487fedB8A0C07']
  /// / )
}

// Check that all values for the struct at this SAID have default values. I.e.
// nothing was changed due to invalid request
export const checkServiceAgreementAbsent = async (
  coordinator: any,
  serviceAgreementID: any
) => {
  const sa = await coordinator.serviceAgreements.call(toHex(serviceAgreementID))
  assertBigNum(sa[0], bigNum(0))
  assertBigNum(sa[1], bigNum(0))
  assertBigNum(sa[2], bigNum(0))
  assert.equal(
    sa[3],
    '0x0000000000000000000000000000000000000000000000000000000000000000'
  )
}

export const newServiceAgreement = async (params: any): Promise<any> => {
  const agreement: any = {}
  params = params || {}
  agreement.payment = params.payment || 1000000000000000000
  agreement.expiration = params.expiration || 300
  agreement.endAt = params.endAt || sixMonthsFromNow()
  agreement.oracles = params.oracles || [oracleNode]
  agreement.oracleSignatures = []
  agreement.requestDigest =
    params.requestDigest ||
    newHash(
      '0xbadc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5'
    )

  const sAID = calculateSAID(agreement)
  agreement.id = toHex(sAID)

  for (let i = 0; i < agreement.oracles.length; i++) {
    const oracle = agreement.oracles[i]
    const oracleSignature = await personalSign(oracle, sAID)
    const requestDigestAddr = recoverPersonalSignature(sAID, oracleSignature)
    assert.equal(oracle.toLowerCase(), toHex(requestDigestAddr))
    agreement.oracleSignatures[i] = oracleSignature
  }
  return agreement
}

export const sixMonthsFromNow = (): number =>
  Math.round(Date.now() / 1000.0) + 6 * 30 * 24 * 60 * 60

export const fulfillOracleRequest = async (
  oracle: any,
  request: any,
  response: any,
  options: any
): Promise<any> => {
  if (!options) {
    options = {}
  }

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

export const cancelOracleRequest = async (
  oracle: any,
  request: any,
  options: any
): Promise<any> => {
  if (!options) {
    options = {}
  }

  return oracle.cancelOracleRequest(
    request.id,
    request.payment,
    request.callbackFunc,
    request.expiration,
    options
  )
}
