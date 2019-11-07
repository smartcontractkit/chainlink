import { ethers } from 'ethers'
import { createFundedWallet } from './wallet'
import { assert } from 'chai'
import { Oracle } from './generated/Oracle'
import { generated as chainlinkv05 } from 'chainlinkv0.5'
import { LinkToken } from './generated/LinkToken'
import { makeDebug } from './debug'
import cbor from 'cbor'
import { EmptyOracle } from './generated/EmptyOracle'
import { OracleFactory } from './generated/OracleFactory'

const debug = makeDebug('helpers')

export interface Roles {
  defaultAccount: ethers.Wallet
  oracleNode: ethers.Wallet
  oracleNode1: ethers.Wallet
  oracleNode2: ethers.Wallet
  oracleNode3: ethers.Wallet
  stranger: ethers.Wallet
  consumer: ethers.Wallet
}

export interface Personas {
  Default: ethers.Wallet
  Neil: ethers.Wallet
  Ned: ethers.Wallet
  Nelly: ethers.Wallet
  Carol: ethers.Wallet
  Eddy: ethers.Wallet
}

interface RolesAndPersonas {
  roles: Roles
  personas: Personas
}

/**
 * This helper function allows us to make use of ganache snapshots,
 * which allows us to snapshot one state instance and revert back to it.
 *
 * This is used to memoize expensive setup calls typically found in beforeEach hooks when we
 * need to setup our state with contract deployments before running assertions.
 *
 * @param provider The provider that's used within the tests
 * @param cb The callback to execute that generates the state we want to snapshot
 */
export function useSnapshot(
  provider: ethers.providers.JsonRpcProvider,
  cb: () => Promise<void>,
) {
  const d = debug.extend('memoizeDeploy')
  let hasDeployed = false
  let snapshotId = ''

  return async () => {
    if (!hasDeployed) {
      d('executing deployment..')
      await cb()

      d('snapshotting...')
      /* eslint-disable-next-line require-atomic-updates */
      snapshotId = await provider.send('evm_snapshot', undefined)
      d('snapshot id:%s', snapshotId)

      /* eslint-disable-next-line require-atomic-updates */
      hasDeployed = true
    } else {
      d('reverting to snapshot: %s', snapshotId)
      await provider.send('evm_revert', snapshotId)

      d('re-creating snapshot..')
      /* eslint-disable-next-line require-atomic-updates */
      snapshotId = await provider.send('evm_snapshot', undefined)
      d('recreated snapshot id:%s', snapshotId)
    }
  }
}

/**
 * A wrapper function to make generated contracts compatible with truffle test suites.
 *
 * Note that the returned contract is an instance of ethers.Contract, not a @truffle/contract, so there are slight
 * api differences, though largely the same.
 *
 * @see https://docs.ethers.io/ethers.js/html/api-contract.html
 * @param contractFactory The ethers based contract factory to interop with
 * @param address The address to supply as the signer
 */
export function create<T extends new (...args: any[]) => any>(
  contractFactory: T,
  address: string,
): InstanceType<T> {
  const web3Instance = (global as any).web3
  const provider = new ethers.providers.Web3Provider(
    web3Instance.currentProvider,
  )
  const signer = provider.getSigner(address)
  const factory = new contractFactory(signer)

  return factory
}

/**
 * Generate roles and personas for tests along with their corrolated account addresses
 */
export async function initializeRolesAndPersonas(
  provider: ethers.providers.JsonRpcProvider,
): Promise<RolesAndPersonas> {
  const accounts = await Promise.all(
    Array(6)
      .fill(null)
      .map(async (_, i) => createFundedWallet(provider, i).then(w => w.wallet)),
  )

  const personas: Personas = {
    Default: accounts[0],
    Neil: accounts[1],
    Ned: accounts[2],
    Nelly: accounts[3],
    Carol: accounts[4],
    Eddy: accounts[5],
  }

  const roles: Roles = {
    defaultAccount: accounts[0],
    oracleNode: accounts[1],
    oracleNode1: accounts[1],
    oracleNode2: accounts[2],
    oracleNode3: accounts[3],
    stranger: accounts[4],
    consumer: accounts[5],
  }

  return { personas, roles }
}

type AsyncFunction = () => Promise<void>
export async function assertActionThrows(action: AsyncFunction) {
  let e: Error | undefined = undefined
  try {
    await action()
  } catch (error) {
    e = error
  }
  if (!e) {
    assert.exists(e, 'Expected an error to be raised')
    return
  }

  assert(e.message, 'Expected an error to contain a message')

  const ERROR_MESSAGES = ['invalid opcode', 'revert']
  const hasErrored = ERROR_MESSAGES.reduce(
    (prev, next) => prev || e!.message.includes(next),
    false,
  )

  assert(
    hasErrored,
    `expected following error message to include ${ERROR_MESSAGES.join(
      ' or ',
    )}. Got: "${e.message}"`,
  )
}

export function checkPublicABI(
  contract: ethers.Contract | ethers.ContractFactory,
  expectedPublic: string[],
) {
  const actualPublic = []
  for (const method of contract.interface.abi) {
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

export const { utils } = ethers
/**
 * Convert a value to a hex string
 * @param args Value to convert to a hex string
 */
export function toHex(
  ...args: Parameters<typeof utils.hexlify>
): ReturnType<typeof utils.hexlify> {
  return utils.hexlify(...args)
}

/**
 * Convert an Ether value to a wei amount
 * @param args Ether value to convert to an Ether amount
 */
export function toWei(
  ...args: Parameters<typeof utils.parseEther>
): ReturnType<typeof utils.parseEther> {
  return utils.parseEther(...args)
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
    requester: requester,
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

/**
 * Decodes a CBOR hex string, and adds opening and closing brackets to the CBOR if they are not present.
 *
 * @param hexstr The hex string to decode
 */
export function decodeDietCBOR(hexstr: string) {
  const buf = hexToBuf(hexstr)

  return cbor.decodeFirstSync(addCBORMapDelimiters(buf))
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

/**
 * Add a starting and closing map characters to a CBOR encoding if they are not already present.
 */
function addCBORMapDelimiters(buffer: Buffer): Buffer {
  if (buffer[0] >> 5 === 5) {
    return buffer
  }

  /**
   * This is the opening character of a CBOR map.
   * @see https://en.wikipedia.org/wiki/CBOR#CBOR_data_item_header
   */
  const startIndefiniteLengthMap = Buffer.from([0xbf])
  /**
   * This is the closing character in a CBOR map.
   * @see https://en.wikipedia.org/wiki/CBOR#CBOR_data_item_header
   */
  const endIndefiniteLengthMap = Buffer.from([0xff])
  return Buffer.concat(
    [startIndefiniteLengthMap, buffer, endIndefiniteLengthMap],
    buffer.length + 2,
  )
}

/**
 * Add a hex prefix to a hex string
 * @param hex The hex string to prepend the hex prefix to
 */
export function addHexPrefix(hex: string): string {
  return hex.startsWith('0x') ? hex : `0x${hex}`
}

export function stripHexPrefix(hex: string): string {
  if (!ethers.utils.isHexString(hex)) {
    throw Error(`Expected valid hex string, got: "${hex}"`)
  }

  return hex.replace('0x', '')
}

/**
 * Convert a number value to bytes32 format
 *
 * @param num The number value to convert to bytes32 format
 */
export function numToBytes32(
  num: Parameters<typeof ethers.utils.hexlify>[0],
): string {
  const hexNum = ethers.utils.hexlify(num)
  const strippedNum = stripHexPrefix(hexNum)
  if (strippedNum.length > 32 * 2) {
    throw Error(
      'Cannot convert number to bytes32 format, value is greater than maximum bytes32 value',
    )
  }
  return addHexPrefix(strippedNum.padStart(32 * 2, '0'))
}

export function toUtf8(
  ...args: Parameters<typeof ethers.utils.toUtf8Bytes>
): ReturnType<typeof ethers.utils.toUtf8Bytes> {
  return ethers.utils.toUtf8Bytes(...args)
}

/**
 * Compute the keccak256 cryptographic hash of a value, returned as a hex string.
 * (Note: often Ethereum documentation refers to this, incorrectly, as SHA3)
 * @param args The data to compute the keccak256 hash of
 */
export function keccak(
  ...args: Parameters<typeof ethers.utils.keccak256>
): ReturnType<typeof ethers.utils.keccak256> {
  return utils.keccak256(...args)
}

type TxOptions = Omit<ethers.providers.TransactionRequest, 'to' | 'from'>

export async function fulfillOracleRequest(
  oracleContract: Oracle | EmptyOracle,
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

export async function cancelOracleRequest(
  oracleContract: Oracle | EmptyOracle,
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

export function requestDataBytes(
  specId: string,
  to: string,
  fHash: string,
  nonce: number,
  dataBytes: string,
): string {
  const ocFactory = new OracleFactory()

  return ocFactory.interface.functions.oracleRequest.encode([
    ethers.constants.AddressZero,
    0,
    specId,
    to,
    fHash,
    nonce,
    1,
    dataBytes,
  ])
}

// link param must be from linkContract(), if amount is a BN
export function requestDataFrom(
  oc: Oracle,
  link: LinkToken,
  amount: ethers.utils.BigNumberish,
  args: string,
  options: Omit<ethers.providers.TransactionRequest, 'to' | 'from'> = {},
): ReturnType<typeof link.transferAndCall> {
  if (!options) {
    options = { value: 0 }
  }

  return link.transferAndCall(oc.address, amount, args, options)
}

export async function increaseTime5Minutes(
  provider: ethers.providers.JsonRpcProvider,
): Promise<void> {
  await provider.send('evm_increaseTime', [300])
}

/**
 * Convert a buffer to a hex string
 * @param hexstr The hex string to convert to a buffer
 */
export function hexToBuf(hexstr: string): Buffer {
  return Buffer.from(stripHexPrefix(hexstr), 'hex')
}

const { CoordinatorFactory } = chainlinkv05
type Hash = ReturnType<typeof ethers.utils.keccak256>
type Coordinator = ReturnType<chainlinkv05.CoordinatorFactory['attach']>
type ServiceAgreement = Parameters<Coordinator['initiateServiceAgreement']>[0]

/**
 * Digest of the ServiceAgreement.
 */
export const generateSAID = (sa: ServiceAgreement): Hash => {
  const [saParam] = new CoordinatorFactory().interface.functions.getId.inputs
  if (
    saParam.name !== '_agreement' ||
    saParam.type !== 'tuple' ||
    !saParam.components
  ) {
    throw Error(
      `extracted wrong version of struct tuple: ${saParam} from coordinatorFactory.interface.functions.getId`,
    )
  }

  const abiValues = saParam.components.reduce(
    (prev, next) => {
      prev.types.push(next.type)
      prev.values.push(sa[next.name as keyof ServiceAgreement])

      return prev
    },
    { types: [], values: [] },
  )

  return ethers.utils.solidityKeccak256(abiValues.types, abiValues.values)
}
