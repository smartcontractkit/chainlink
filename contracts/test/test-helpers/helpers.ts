import { BigNumber, BigNumberish, Contract, ContractTransaction } from 'ethers'
import { providers } from 'ethers'
import { assert, expect } from 'chai'
import hre, { ethers, network } from 'hardhat'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import cbor from 'cbor'
import { LinkToken } from '../../typechain'

/**
 * Convert string to hex bytes
 * @param data string to convert to hex bytes
 */
export function stringToBytes(data: string): string {
  return ethers.utils.hexlify(ethers.utils.toUtf8Bytes(data))
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

/**
 * Retrieve single log from transaction
 *
 * @param tx The transaction to wait for, then extract logs from
 * @param index The index of the log to retrieve
 */
export async function getLog(
  tx: ContractTransaction,
  index: number,
): Promise<providers.Log> {
  const logs = await getLogs(tx)
  if (!logs[index]) {
    throw Error('unable to extract log from transaction receipt')
  }
  return logs[index]
}

/**
 * Extract array of logs from a transaction
 *
 * @param tx The transaction to wait for, then extract logs from
 */
export async function getLogs(
  tx: ContractTransaction,
): Promise<providers.Log[]> {
  const receipt = await tx.wait()
  if (!receipt.logs) {
    throw Error('unable to extract logs from transaction receipt')
  }
  return receipt.logs
}

/**
 * Convert a UTF-8 string into a bytes32 hex string representation
 *
 * The inverse function of [[parseBytes32String]]
 *
 * @param args The UTF-8 string representation to convert to a bytes32 hex string representation
 */
export function toBytes32String(
  ...args: Parameters<typeof ethers.utils.formatBytes32String>
): ReturnType<typeof ethers.utils.formatBytes32String> {
  return ethers.utils.formatBytes32String(...args)
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
 * Create a buffer from a hex string
 *
 * @param hexstr The hex string to convert to a buffer
 */
export function hexToBuf(hexstr: string): Buffer {
  return Buffer.from(stripHexPrefix(hexstr), 'hex')
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
 * Convert an Ether value to a wei amount
 *
 * @param args Ether value to convert to an Ether amount
 */
export function toWei(
  ...args: Parameters<typeof ethers.utils.parseEther>
): ReturnType<typeof ethers.utils.parseEther> {
  return ethers.utils.parseEther(...args)
}

/**
 * Converts any number, BigNumber, hex string or Arrayish to a hex string.
 *
 * @param args Value to convert to a hex string
 */
export function toHex(
  ...args: Parameters<typeof ethers.utils.hexlify>
): ReturnType<typeof ethers.utils.hexlify> {
  return ethers.utils.hexlify(...args)
}

/**
 * Increase the current time within the evm to 5 minutes past the current time
 *
 * @param provider The ethers provider to send the time increase request to
 */
export async function increaseTime5Minutes(
  provider: providers.JsonRpcProvider,
): Promise<void> {
  await increaseTimeBy(5 * 60, provider)
}

/**
 * Increase the current time within the evm to "n" seconds past the current time
 *
 * @param seconds The number of seconds to increase to the current time by
 * @param provider The ethers provider to send the time increase request to
 */
export async function increaseTimeBy(
  seconds: number,
  provider: providers.JsonRpcProvider,
) {
  await provider.send('evm_increaseTime', [seconds])
}

/**
 * Instruct the provider to mine an additional block
 *
 * @param provider The ethers provider to instruct to mine an additional block
 */
export async function mineBlock(provider: providers.JsonRpcProvider) {
  await provider.send('evm_mine', [])
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
  return ethers.utils.getAddress(hex.slice(26))
}

/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
export function publicAbi(contract: Contract, expectedPublic: string[]) {
  const actualPublic = []
  for (const m in contract.functions) {
    if (!m.includes('(')) {
      actualPublic.push(m)
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

/**
 * Converts an L1 address to an Arbitrum L2 address
 *
 * @param l1Address Address on L1
 */
export function toArbitrumL2AliasAddress(l1Address: string): string {
  return ethers.utils.getAddress(
    BigNumber.from(l1Address)
      .add('0x1111000000000000000000000000000000001111')
      .toHexString()
      .replace('0x01', '0x'),
  )
}

/**
 * Lets you impersonate and sign transactions from any account.
 *
 * @param address Address to impersonate
 */
export async function impersonateAs(
  address: string,
): Promise<SignerWithAddress> {
  await hre.network.provider.request({
    method: 'hardhat_impersonateAccount',
    params: [address],
  })
  return await ethers.getSigner(address)
}

export async function stopImpersonateAs(address: string): Promise<void> {
  await hre.network.provider.request({
    method: 'hardhat_stopImpersonatingAccount',
    params: [address],
  })
}

export async function assertBalance(
  address: string,
  balance: BigNumberish,
  msg?: string,
) {
  expect(await ethers.provider.getBalance(address)).equal(balance, msg)
}

export async function assertLinkTokenBalance(
  lt: LinkToken,
  address: string,
  balance: BigNumberish,
  msg?: string,
) {
  expect(await lt.balanceOf(address)).equal(balance, msg)
}

export async function assertSubscriptionBalance(
  coordinator: Contract,
  subID: BigNumberish,
  balance: BigNumberish,
  msg?: string,
) {
  expect((await coordinator.getSubscription(subID)).balance).deep.equal(
    balance,
    msg,
  )
}

export async function setTimestamp(timestamp: number) {
  await network.provider.request({
    method: 'evm_setNextBlockTimestamp',
    params: [timestamp],
  })
  await network.provider.request({
    method: 'evm_mine',
    params: [],
  })
}

export async function fastForward(duration: number) {
  await network.provider.request({
    method: 'evm_increaseTime',
    params: [duration],
  })
  await network.provider.request({
    method: 'evm_mine',
    params: [],
  })
}

export async function reset() {
  await network.provider.request({
    method: 'hardhat_reset',
    params: [],
  })
}
