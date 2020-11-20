/**
 * @packageDocumentation
 *
 * This file provides common utility functions to perform ethereum related tasks, like
 * data format manipulation of buffers and hex strings,
 * finding, accessing logs and events,
 * and increasing test evm time.
 */
import cbor from 'cbor'
import { assert } from 'chai'
import { ethers, utils } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { EventDescription } from 'ethers/utils'

/**
 * Convert string to hex bytes
 * @param data string to onvert to hex bytes
 */
export function stringToBytes(data: string): string {
  return ethers.utils.hexlify(ethers.utils.toUtf8Bytes(data))
}

/**
 * Convert hex bytes to utf8 string
 * @param data bytes to convert to utf8 stirng
 */
export function bytesToString(data: string): string {
  return ethers.utils.toUtf8String(data)
}

/**
 * Parse out an evm word (32 bytes) into an address (20 bytes) representation
 *
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
 * Convert a number value to bytes32 format
 *
 * @param num The number value to convert to bytes32 format
 */
export function numToBytes32(num: Parameters<typeof utils.hexlify>[0]): string {
  const hexNum = utils.hexlify(num)
  const strippedNum = stripHexPrefix(hexNum)
  if (strippedNum.length > 32 * 2) {
    throw Error(
      'Cannot convert number to bytes32 format, value is greater than maximum bytes32 value',
    )
  }
  return addHexPrefix(strippedNum.padStart(32 * 2, '0'))
}

/**
 * Convert a UTF-8 string into a bytes32 hex string representation
 *
 * The inverse function of [[parseBytes32String]]
 *
 * @param args The UTF-8 string representation to convert to a bytes32 hex string representation
 */
export function toBytes32String(
  ...args: Parameters<typeof utils.formatBytes32String>
): ReturnType<typeof utils.formatBytes32String> {
  return utils.formatBytes32String(...args)
}

/**
 * Convert a bytes32 formatted hex string into its UTF-8 representation
 *
 * The inverse function of [[toBytes32String]].
 *
 * @param args The bytes32 hex string representation to convert to an UTF-8 representation
 */
export function parseBytes32String(
  ...args: Parameters<typeof utils.parseBytes32String>
): ReturnType<typeof utils.parseBytes32String> {
  return utils.parseBytes32String(...args)
}

/**
 * Converts any number, BigNumber, hex string or Arrayish to a hex string.
 *
 * @param args Value to convert to a hex string
 */
export function toHex(
  ...args: Parameters<typeof utils.hexlify>
): ReturnType<typeof utils.hexlify> {
  return utils.hexlify(...args)
}

/**
 * Create a buffer from a hex string
 *
 * @param hexstr The hex string to convert to a buffer
 */
export function hexToBuf(hexstr: string): Buffer {
  return Buffer.from(stripHexPrefix(hexstr), 'hex')
}

/**
 * Convert an Ether value to a wei amount
 *
 * @param args Ether value to convert to an Ether amount
 */
export function toWei(
  ...args: Parameters<typeof utils.parseEther>
): ReturnType<typeof utils.parseEther> {
  return utils.parseEther(...args)
}

/**
 * Convert a value to an ethers BigNum
 *
 * @param num Value to convert to a BigNum
 */
export function bigNum(num: utils.BigNumberish): utils.BigNumber {
  return utils.bigNumberify(num)
}

/**
 * Convert a UTF-8 string into a bytearray
 *
 * @param args The values needed to convert a string into a bytearray
 */
export function toUtf8Bytes(
  ...args: Parameters<typeof utils.toUtf8Bytes>
): ReturnType<typeof utils.toUtf8Bytes> {
  return utils.toUtf8Bytes(...args)
}

/**
 * Turn a [x,y] coordinate into an ethereum address
 *
 * @param pubkey The x,y coordinate to turn into an ethereum address
 */
export function pubkeyToAddress(pubkey: utils.BigNumber[]) {
  // transform the value according to what ethers expects as a value
  const concatResult = `0x04${pubkey
    .map((coord) => coord.toHexString())
    .join('')
    .replace(/0x/gi, '')}`

  return utils.computeAddress(concatResult)
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
 *
 * @param hex The hex string to prepend the hex prefix to
 */
export function addHexPrefix(hex: string): string {
  return hex.startsWith('0x') ? hex : `0x${hex}`
}

/**
 * Strip the leading 0x hex prefix from a hex string
 *
 * @param hex The hex string to strip the leading hex prefix out of
 */
export function stripHexPrefix(hex: string): string {
  if (!ethers.utils.isHexString(hex)) {
    throw Error(`Expected valid hex string, got: "${hex}"`)
  }

  return hex.replace('0x', '')
}

/**
 * Compute the keccak256 cryptographic hash of a value, returned as a hex string.
 * (Note: often Ethereum documentation refers to this, incorrectly, as SHA3)
 *
 * @param args The data to compute the keccak256 hash of
 */
export function keccak(
  ...args: Parameters<typeof utils.keccak256>
): ReturnType<typeof utils.keccak256> {
  return utils.keccak256(...args)
}

/**
 * Increase the current time within the evm to "n" seconds past the current time
 *
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
 * Instruct the provider to mine an additional block
 *
 * @param provider The ethers provider to instruct to mine an additional block
 */
export async function mineBlock(provider: ethers.providers.JsonRpcProvider) {
  await provider.send('evm_mine', [])
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

/**
 * Find an event within a transaction receipt by its event description
 *
 * @param receipt The events array to search through
 * @param eventDescription The event description to pass to check its name by
 */
export function findEventIn(
  receipt: ContractReceipt,
  eventDescription: EventDescription,
): ethers.Event | undefined {
  // the first topic of a log is always the keccak-256 hash of the event signature
  const event = receipt.events?.find(
    (e) => e.topics[0] === eventDescription.topic,
  )

  return event
}

/**
 * Calculate six months from the current date in seconds
 */
export function sixMonthsFromNow(): utils.BigNumber {
  return utils.bigNumberify(
    Math.round(Date.now() / 1000.0) + 6 * 30 * 24 * 60 * 60,
  )
}

/**
 * Extract array of logs from a transaction
 *
 * @param tx The transaction to wait for, then extract logs from
 */
export async function getLogs(
  tx: ethers.ContractTransaction,
): Promise<ethers.providers.Log[]> {
  const receipt = await tx.wait()
  if (!receipt.logs) {
    throw Error('unable to extract logs from transaction receipt')
  }
  return receipt.logs
}

/**
 * Retrieve single log from transaction
 *
 * @param tx The transaction to wait for, then extract logs from
 * @param index The index of the log to retrieve
 */
export async function getLog(
  tx: ethers.ContractTransaction,
  index: number,
): Promise<ethers.providers.Log> {
  const logs = await getLogs(tx)
  if (!logs[index]) {
    throw Error('unable to extract log from transaction receipt')
  }
  return logs[index]
}
