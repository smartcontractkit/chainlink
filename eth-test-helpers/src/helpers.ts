import cbor from 'cbor'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { EventDescription } from 'ethers/utils'

export const { utils } = ethers

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
 * Parse out an evm word (32 bytes) into an address (20 bytes) representation
 * @param hex The evm word in hex string format to parse the address
 * out of.
 */
export function evmWordToAddress(hex?: string): string {
  if (!hex) {
    throw Error('Input not defined')
  }

  assert.equal(hex.slice(0, 26), '0x000000000000000000000000')
  return utils.getAddress(hex.slice(26))
}

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

/**
 * Decodes a CBOR hex string, and adds opening and closing brackets to the CBOR if they are not present.
 *
 * @param hexstr The hex string to decode
 */
export function decodeDietCBOR(hexstr: string) {
  const buf = hexToBuf(hexstr)

  return cbor.decodeFirstSync(addCBORMapDelimiters(buf))
}

/**
 * Add a starting and closing map characters to a CBOR encoding if they are not already present.
 */
export function addCBORMapDelimiters(buffer: Buffer): Buffer {
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

/**
 * Increase the current time within the evm to "n" seconds past the current time
 * @param seconds The number of seconds to increase to the current time by
 * @param provider The ethers provider to send the time increase request to
 */
export async function increaseTimeBy(
  seconds: number,
  provider: ethers.providers.JsonRpcProvider,
) {
  await provider.send('evm_increaseTime', [seconds])
}

/**
 * Increase the current time within the evm to 5 minutes past the current time
 *
 * @param provider The ethers provider to send the time increase request to
 */
export async function increaseTime5Minutes(
  provider: ethers.providers.JsonRpcProvider,
): Promise<void> {
  await increaseTimeBy(5 * 600, provider)
}

/**
 * Convert a buffer to a hex string
 * @param hexstr The hex string to convert to a buffer
 */
export function hexToBuf(hexstr: string): Buffer {
  return Buffer.from(stripHexPrefix(hexstr), 'hex')
}

/**
 * Turn a [x,y] coordinate into an ethereum address
 * @param pubkey The x,y coordinate to turn into an ethereum address
 */
export function pubkeyToAddress(pubkey: ethers.utils.BigNumber[]) {
  // transform the value according to what ethers expects as a value
  const concatResult = `0x04${pubkey
    .map(coord => coord.toHexString())
    .join('')
    .replace(/0x/gi, '')}`

  return ethers.utils.computeAddress(concatResult)
}

interface EventArgsArray extends Array<any> {
  [key: string]: any
}
/**
 * Typecast an ethers event to its proper type, until
 * https://github.com/ethers-io/ethers.js/pull/698 is addressed
 *
 * @param event The event to typecast
 */
export function eventArgs(event?: ethers.Event) {
  return (event?.args as any) as EventArgsArray
}

export interface TypedEventDescription<
  T extends Pick<EventDescription, 'encodeTopics'>
> extends EventDescription {
  encodeTopics: T['encodeTopics']
}

/**
 * Find an event within a transaction receipt by its event description
 *
 * @param receipt The events array to search through
 * @param eventDescription The event description to pass to check its name by
 */
export function findEventIn(
  receipt: ContractReceipt,
  eventDescription: TypedEventDescription<any>,
): ethers.Event | undefined {
  return receipt.events?.find(e => e.event === eventDescription.name)
}

/**
 * Calculate six months from the current date in seconds
 */
export function sixMonthsFromNow(): ethers.utils.BigNumber {
  return ethers.utils.bigNumberify(
    Math.round(Date.now() / 1000.0) + 6 * 30 * 24 * 60 * 60,
  )
}
