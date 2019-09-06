import { BN } from 'bn.js'
import cbor from 'cbor'
import * as abi from 'ethereumjs-abi'
import * as util from 'ethereumjs-util'
import { FunctionFragment, ParamType } from 'ethers/utils/abi-coder'
import TruffleContract from 'truffle-contract'
import { linkToken } from './linkToken'
import { assertBigNum } from './matchers'

const HEX_BASE = 16

// https://github.com/ethereum/web3.js/issues/1119#issuecomment-394217563
web3.providers.HttpProvider.prototype.sendAsync =
  web3.providers.HttpProvider.prototype.send
export const eth = web3.eth

const INVALIDVALUE: Record<string, any> = {
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
export let personas: Record<string, any> = {}

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

  // allow personas instead of roles
  personas.Default = defaultAccount
  personas.Neil = oracleNode1
  personas.Ned = oracleNode2
  personas.Nelly = oracleNode3
  personas.Carol = consumer
  personas.Eddy = stranger
})

const bNToStringOrIdentity = (a: any): any => (BN.isBN(a) ? a.toString() : a)
export const BNtoUint8Array = (n: BN): Uint8Array =>
  Uint8Array.from((new BN(n)).toArray('be', 32))

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

export const linkContract = async (account: any): Promise<any> => {
  account = account || defaultAccount
  const receipt = await web3.eth.sendTransaction({
    data: linkToken.bytecode,
    from: account,
    gasLimit: 2000000
  })
  const contract = TruffleContract({ abi: linkToken.abi })
  contract.setProvider(web3.currentProvider)
  contract.defaults({
    from: account,
    gas: 3500000,
    gasPrice: 10000000000
  })

  return wrappedERC20(await contract.at(receipt.contractAddress))
}

export const bigNum = web3.utils.toBN
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

const hexRegExp = /^(0[xX])?[0-9a-fA-F]+$/
const isHex = hexRegExp.test.bind(hexRegExp)

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
    assert.equal(typeof arg, 'string', `Don't know how to convert ${arg} to hex`)
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

export const getEvents = (contract: any): Promise<any[]> =>
  new Promise((resolve, reject) =>
    contract
      .getPastEvents('allEvents', { fromBlock: 1 }) // https://ethereum.stackexchange.com/questions/71307/mycontract-getpasteventsallevents-returns-empty-array
      .then((events: any) => resolve(events))
      .catch((error: any) => reject(error))
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
    options = { value: 0 }
  }
  return link.transferAndCall(oc.address, amount, args, options)
}

// ABI specification for the given method on the given contract
export const getMethod = (contract: TruffleContract, methodName: string): FunctionFragment => {
  const methodABIs = contract.abi.filter(
    ({ name: attrName }: FunctionFragment) => attrName == methodName
  )
  const fqName = `${contract.contractName}.${methodName}: ${methodABIs}`
  assert.equal(methodABIs.length, 1, `No method ${fqName}, or ambiguous`)
  return methodABIs[0]
}

export const functionSelector = web3.eth.abi.encodeFunctionSignature

export const functionSelectorFromAbi = (contract: TruffleContract, name: string): string => 
  functionSelector(getMethod(contract, name))

export const assertActionThrows = (action: any, messageContains?: RegExp) =>
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
      assert(
        (!messageContains) || messageContains.test(errorMessage),
        `expected error message to contain ${messageContains}: ${errorMessage}`
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
    callbackFunc,
    data: autoAddMapDelimiters(data),
    dataVersion: version,
    expiration,
    id: requestId,
    jobId: log.topics[1],
    payment,
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

// newUint8ArrayFromHex returns count bytes from hex string. They are padded at
// the end to have `count` bytes
export const newUint8ArrayFromHex = (str: string, count: number): Uint8Array => {
  assert(isHex(str), `${str} is not a hexadecimal value`)
  const hexCount = 2 * count + 2
  const hexStr = Ox(str)
  assert(hexStr.length <= hexCount, `${str} won't fit in ${count} bytes`)
  return Uint8Array.from(web3.utils.hexToBytes(hexStr.padEnd(hexCount, '0')))
}

// newUint8ArrayFromDecimal returns count bytes from a decimal number. They are
// padded at the start to have `count` bytes
export const newUint8ArrayFromDecimal = (str: string, count: number): any => {
  assert(/^[0-9]+$/.test(str), `${str} is not a decimal value.`)
  const hexCount = 2 * count
  const hexStr = (new BN(str, 10)).toString(16).padStart(hexCount, '0')
  assert(hexStr.length <= hexCount, `${str} won't fit in ${count} bytes`)
  return Uint8Array.from(web3.utils.hexToBytes(Ox(hexStr)))
}

// newSignature returns a signature object with v, r, and s broken up
export const newSignature = (str: string): any => {
  const oracleSignature = newUint8ArrayFromHex(str, 65)
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

/**
 * @param str hexadecimal/decimal representation of integer. 0x must prefix hex.
 * @returns 32-byte representation of number. If hexadecimal, zero-padded on
 *          the right. If decimal, zero-padded on the left.
 * @todo (alx): Split this into more specific and explicit functions.
 */
export const newHash = (str: string): Uint8Array => 
  (/^0[xX]/.test(str) ? newUint8ArrayFromHex : newUint8ArrayFromDecimal)(
  str,
  32
)

// newAddress returns a 20-byte Uint8Array for representing an address
export const newAddress = (str: string): Uint8Array => {
  return newUint8ArrayFromHex(str, 20)
}

export const newSelector = (str: string): Uint8Array =>
  newUint8ArrayFromHex(str, 4)

export const toBuffer = (uint8a: Uint8Array): Buffer => Buffer.from(uint8a)

export const concatUint8Arrays = (...arrays: Uint8Array[]): Uint8Array =>
  Uint8Array.from(Buffer.concat(arrays.map(Buffer.from)))

export const increaseTime5Minutes = async () => {
  await web3.currentProvider.send(
    {
      id: 0,
      jsonrpc: '2.0',
      method: 'evm_increaseTime',
      params: [300]
    },
    (error: any) => {
      if (error) {
        throw Error(`Error during helpers.increaseTime5Minutes! ${error}`)
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
    (error: any) => {
      if (error) {
        throw Error(`Error during ${evmMethod}! ${error}`)
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

interface Signature {v: number, r: Uint8Array, s: Uint8Array}

interface ServiceAgreement { // Corresponds to ServiceAgreement struct in CoordinatorInterface.sol
  payment: BN, // uint256
  expiration: BN, // uint256
  endAt: BN, // uint256
  oracles: string[], // 0x hex representation of oracle addresses (uint160's)
  requestDigest: string, // 0x hex representation of bytes32
  aggregator: string, // 0x hex representation of aggregator address
  aggInitiateJobSelector: string, // 0x hex representation of aggregator.initiateAggregatorForJob function selector (uint32)
  aggFulfillSelector: string, // function selector for aggregator.fulfill
  // Information which is useful to carry around with the agreement, but not
  // part of the solidity struct
  id: string // ServiceAgreement Id (sAId)
  oracleSignatures: Signature[]
}

export const calculateSAID = (sa: ServiceAgreement): Uint8Array => {
  const serviceAgreementIDInput = concatUint8Arrays(
    BNtoUint8Array(sa.payment),
    BNtoUint8Array(sa.expiration),
    BNtoUint8Array(sa.endAt),
    // Each address in this list is padded to a uint256, despite being a uint160
    ...sa.oracles.map(pad0xHexTo256Bit).map(newHash),
    newHash(sa.requestDigest),
    newAddress(sa.aggregator),
    newSelector(sa.aggInitiateJobSelector),
    newSelector(sa.aggFulfillSelector)
  )
  const serviceAgreementIDInputDigest = util.keccak(
    toHex(serviceAgreementIDInput)
  )
  return newHash(toHex(serviceAgreementIDInputDigest))
}

export const recoverPersonalSignature = (
  message: Uint8Array,
  signature: Signature
): any => {
  const personalSignPrefix = newUint8ArrayFromStr(
    '\x19Ethereum Signed Message:\n'
  )
  const personalSignMessage = Uint8Array.from(
    concatUint8Arrays(
      personalSignPrefix,
      newUint8ArrayFromStr(message.length.toString()),
      message
    )
  )
  const digest = util.keccak(toBuffer(personalSignMessage))
  const requestDigestPubKey = util.ecrecover(
    digest,
    signature.v,
    Buffer.from(signature.r),
    Buffer.from(signature.s)
  )
  return util.pubToAddress(requestDigestPubKey)
}

export const personalSign = async (
  account: any,
  message: any
): Promise<any> => {
  const eMsg = `Message ${message} is not a recognized representation of a ` +
    'byte array. (Can be Buffer, BigNumber, Uint8Array, 0x-prepended ' +
    'hexadecimal string.)'
  assert(isByteRepresentation(message), eMsg)
  return newSignature(await web3.eth.sign(toHex(message), account))
}

export const executeServiceAgreementBytes = (
  sAID: any,
  callbackAddr: any,
  callbackFunctionId: any,
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
  const values = [0, 0, sAID, callbackAddr, callbackFunctionId, nonce, 1, data]
  const encoded = abiEncode(types, values)
  const funcSelector = functionSelector(
    'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)'
  )
  return funcSelector + encoded
}

// Convenience functions for constructing hexadecimal representations of
// binary serializations.export const padHexTo256Bit = (s: string): string => s.padStart(64, '0')
export const strip0x = (s: string): string =>
  /^0[xX]/.test(s) ? s.slice(2) : s
export const padHexTo256Bit = (s: string): string => s.padStart(64, '0')
export const pad0xHexTo256Bit = (s: string): string =>
  Ox(padHexTo256Bit(strip0x(s)))
export const padNumTo256Bit = (n: number): string =>
  padHexTo256Bit(n.toString(16))

export const constructStructArgs = (
  fieldNames: string[],
  values: any[]
): any[] => {
  assert.equal(fieldNames.length, values.length)
  const args = []
  for (let i = 0; i < fieldNames.length; i++) {
    args[i] = values[i]
    args[fieldNames[i] as any] = values[i]
  }
  return args
}

// ABI specification for the given argument of the given contract method
const getMethodArg = (
  contract: any,
  methodName: string,
  argName: string
): ParamType => {
  const fqName = `${contract.contractName}.${methodName}`
  const methodABI = getMethod(contract, methodName)
  let eMsg = `${fqName} is not a method: ${methodABI}`
  assert.equal(methodABI.type, 'function', eMsg)
  const argMatches = methodABI.inputs.filter((a: any) => a.name == argName)
  eMsg = `${fqName} has no argument ${argName}, or name is ambiguous`
  assert.equal(argMatches.length, 1, eMsg)
  return argMatches[0]
}

// Struct as mapping => tuple representation of struct, for use in truffle call
//
// TODO(alx): This does not deal with nested structs. It may be possible to do
// that by making an AbiCoder with a custom CoerceFunc which, given a tuple
// type, checks whether the input value is a map or a sequence, and if a map,
// converts it to a sequence as I'm doing here.
export const structAsTuple = (
  struct: { [fieldName: string]: any},
  contract: TruffleContract,
  methodName: string,
  argName: string
): { abi: ParamType, struct: ArrayLike<any> } => {
  const abi: ParamType = getMethodArg(contract, methodName, argName)
  const eMsg = `${contract.contractName}.${methodName}'s argument ${argName} ` +
    `is not a struct: ${abi}`
  assert.equal(abi.type, 'tuple', eMsg)
  return { abi, struct: abi.components.map(({ name }) => struct[name]) }
}

export const initiateServiceAgreementArgs = (
  coordinator: TruffleContract,
  serviceAgreement: ServiceAgreement
): any[] => {
  const signatures = {
    vs: serviceAgreement.oracleSignatures.map(os => os.v),
    rs: serviceAgreement.oracleSignatures.map(os => os.r),
    ss: serviceAgreement.oracleSignatures.map(os => os.s)
  }
  const tup = (s: any, n: any) =>
    structAsTuple(s, coordinator, 'initiateServiceAgreement', n).struct
  return [tup(serviceAgreement, '_agreement'), tup(signatures, '_signatures')]
}

// Call coordinator contract to initiate the specified service agreement, and
// get the return value
export const initiateServiceAgreementCall = async (
  coordinator: TruffleContract,
  serviceAgreement: ServiceAgreement
) => await coordinator.initiateServiceAgreement.call(
  ...initiateServiceAgreementArgs(coordinator, serviceAgreement)
)

/** Call coordinator contract to initiate the specified service agreement. */
export const initiateServiceAgreement = async (
  coordinator: TruffleContract,
  serviceAgreement: ServiceAgreement
) => coordinator.initiateServiceAgreement(
  ...initiateServiceAgreementArgs(coordinator, serviceAgreement))

/** Check that the given service agreement was stored at the correct location */
export const checkServiceAgreementPresent = async (
  coordinator: TruffleContract,
  serviceAgreement: ServiceAgreement
) => {
  const sa = await coordinator.serviceAgreements.call(serviceAgreement.id)
  assertBigNum(sa[0], bigNum(serviceAgreement.payment), 'expected payment')
  assertBigNum(sa[1], bigNum(serviceAgreement.expiration), 'expected expiration')
  assertBigNum(sa[2], bigNum(serviceAgreement.endAt), 'expected endAt date')
  assert.equal(sa[3], serviceAgreement.requestDigest, 'expected requestDigest')
}

// Check that all values for the struct at this SAID have default values. I.e.
// nothing was changed due to invalid request
export const checkServiceAgreementAbsent = async (
  coordinator: any,
  serviceAgreementID: any
) => {
  const sa = await coordinator.serviceAgreements.call(
    toHex(serviceAgreementID).slice(0, 66)
  )
  assertBigNum(sa[0], bigNum(0), 'service agreement is not absent')
  assertBigNum(sa[1], bigNum(0), 'service agreement is not absent')
  assertBigNum(sa[2], bigNum(0), 'service agreement is not absent')
  assert.equal(
    sa[3],
    '0x0000000000000000000000000000000000000000000000000000000000000000'
  )
}

export const newServiceAgreement = async (
  params: Partial<ServiceAgreement>
): Promise<ServiceAgreement> => {
  const agreement: Partial<ServiceAgreement> = {}
  params = params || {}
  agreement.payment = params.payment || new BN('1000000000000000000', 10)
  agreement.expiration = params.expiration || new BN(300)
  agreement.endAt = params.endAt || sixMonthsFromNow()
  agreement.oracles = params.oracles || [oracleNode]
  agreement.oracleSignatures = []
  agreement.requestDigest = params.requestDigest ||
    '0xbadc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5'
  agreement.aggregator = params.aggregator || '0x3141592653589793238462643383279502884197'
  agreement.aggInitiateJobSelector = params.aggInitiateJobSelector || '0x12345678'
  agreement.aggFulfillSelector = params.aggFulfillSelector || '0x87654321'
  const sAID = calculateSAID(agreement as ServiceAgreement)
  agreement.id = toHex(sAID)

  for (let i = 0; i < agreement.oracles.length; i++) {
    const oracle = agreement.oracles[i]
    const oracleSignature = await personalSign(oracle, sAID)
    const requestDigestAddr = recoverPersonalSignature(sAID, oracleSignature)
    assert.equal(oracle.toLowerCase(), toHex(requestDigestAddr))
    agreement.oracleSignatures[i] = oracleSignature
  }
  return agreement as ServiceAgreement
}

export const sixMonthsFromNow = (): number =>
  new BN(Math.round(Date.now() / 1000.0) + 6 * 30 * 24 * 60 * 60)

export const fulfillOracleRequest = async (
  oracle: any,
  request: any,
  response: any,
  options: any
): Promise<any> => {
  if (!options) {
    options = { value: 0 }
  }

  return oracle.fulfillOracleRequest(
    request.id,
    request.payment,
    request.callbackAddr,
    request.callbackFunc,
    request.expiration,
    toHex(response),
    options
  )
}

export const cancelOracleRequest = async (
  oracle: any,
  request: any,
  options: any
): Promise<any> => {
  if (!options) {
    options = { value: 0 }
  }

  return oracle.cancelOracleRequest(
    request.id,
    request.payment,
    request.callbackFunc,
    request.expiration,
    options
  )
}

type numeric = number | BN

export const hexPadUint256 = (n: BN): string => n.toJSON().padStart(64, '0')
export const encodeUint256 = (int: numeric): string => hexPadUint256(new BN(int))
export const encodeInt256 = (int: numeric): string => hexPadUint256(
  (new BN(int)).toTwos(256)
)
export const encodeAddress = (a: string): string => {
  assert(Ox(a).length <= 40, `${a} is too long to be an address`)
  return Ox(strip0x(a).padStart(40, '0'))
}
